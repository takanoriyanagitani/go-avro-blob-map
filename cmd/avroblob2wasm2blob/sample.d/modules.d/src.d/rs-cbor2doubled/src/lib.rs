use std::sync::RwLock;
use std::collections::BTreeMap;

use ciborium::cbor;
use ciborium::Value;
use ciborium_io::Write;

struct Buffer {
    raw: [u8; 65536],
    size: u16,
}

impl Buffer {
    pub fn clear(&mut self) {
        self.size = 0
    }
}

impl Write for Buffer {
    type Error = &'static str;

    fn write_all(&mut self, data: &[u8]) -> Result<(), Self::Error> {
        let bsz: usize = self.size as usize;
        let before: &mut [u8] = &mut self.raw[bsz..];
        let left: usize = 65536 - bsz;

        let sz: usize = data.len().min(left);
        let source: &[u8] = &data[..sz];
        let target: &mut [u8] = &mut before[..sz];
        target.copy_from_slice(source);
        self.size += sz as u16;
        Ok(())
    }

    fn flush(&mut self) -> Result<(), Self::Error> {
        Ok(())
    }
}

static INPUT: RwLock<Buffer> = RwLock::new(Buffer {
    raw: [0; 65536],
    size: 0,
});

static OUTPUT: RwLock<Buffer> = RwLock::new(Buffer {
    raw: [0; 65536],
    size: 0,
});

fn input_size() -> Result<u16, &'static str> {
    let guard = INPUT.try_read().map_err(|_| "unable to read lock")?;
    let i: &Buffer = &guard;
    Ok(i.size)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn b2b_input_size() -> i32 {
    input_size().ok().map(|u| u.into()).unwrap_or(-1)
}

pub fn input_offset() -> Result<*mut u8, &'static str> {
    let mut guard = INPUT.try_write().map_err(|_| "unable to write lock")?;
    let i: &mut Buffer = &mut guard;
    Ok(i.raw.as_mut_ptr())
}

#[allow(unsafe_code)]
#[no_mangle]
pub fn b2b_offset_i() -> *mut u8 {
    input_offset().ok().unwrap_or_else(std::ptr::null_mut)
}

fn allocate(sz: u16) -> Result<(), &'static str> {
    let mut guard = INPUT.try_write().map_err(|_| "unable to write lock")?;
    let mi: &mut Buffer = &mut guard;
    mi.size = sz;
    Ok(())
}

#[allow(unsafe_code)]
#[no_mangle]
pub fn b2b_allocate(sz: u16) -> i32 {
    allocate(sz).ok().map(|_| 65536).unwrap_or(-1)
}

fn output_size() -> Result<u16, &'static str> {
    let guard = OUTPUT.try_read().map_err(|_| "unable to read lock")?;
    let i: &Buffer = &guard;
    Ok(i.size)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn b2b_output_size() -> i32 {
    output_size().ok().map(|u| u.into()).unwrap_or(-1)
}

pub fn output_offset() -> Result<*const u8, &'static str> {
    let guard = OUTPUT.try_read().map_err(|_| "unable to read lock")?;
    let i: &Buffer = &guard;
    Ok(i.raw.as_ptr())
}

#[allow(unsafe_code)]
#[no_mangle]
pub fn b2b_offset_o() -> *const u8 {
    output_offset().ok().unwrap_or_else(std::ptr::null)
}

fn input2output(input: &[u8], output: &mut Buffer) -> Result<(), &'static str> {
    let parsed: Value = ciborium::from_reader(input).map_err(|_| "unable to parse")?;
	let mval: Vec<(Value, Value)> = parsed.into_map().map_err(|_| "not a map")?;
	let pairs = mval.into_iter().map(|pair| {
		let (key, val) = pair;
		key.into_text().map(|s| (s, val)).map_err(|_| "not a text")
	});
	let srval: Result<Vec<(String, Value)>, _> = pairs.collect();
	let sval: Vec<(String, Value)> = srval?;
	let mut bval: BTreeMap<String, Value> = BTreeMap::from_iter(sval);
	let fval: Value = bval.remove("f").unwrap_or(Value::Float(0.0));
	let f: f64 = fval.into_float().map_err(|_| "not a float")?;

    let val: Value = cbor!(f*2.0).map_err(|_| "unable to create a value")?;
    ciborium::into_writer(&val, output).map_err(|_| "unable to write")?;
    Ok(())
}

fn convert() -> Result<(), &'static str> {
    let guard = INPUT.try_read().map_err(|_| "unable to read lock")?;
    let input: &Buffer = &guard;
    let inputsz: usize = input.size as usize;
    let limited: &[u8] = &input.raw[..inputsz];

    let mut out = OUTPUT.try_write().map_err(|_| "unable to write lock")?;
    let output: &mut Buffer = &mut out;
    output.clear();
    input2output(limited, output)?;
    Ok(())
}

#[allow(unsafe_code)]
#[no_mangle]
pub fn b2b_convert() -> i32 {
    convert()
        .and_then(|_| output_size())
        .ok()
        .map(|u| u.into())
        .unwrap_or(-1)
}
