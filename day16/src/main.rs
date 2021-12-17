use std::collections::VecDeque;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

struct BitReader {
    bit: usize,
    offset: usize,
    buffer: VecDeque<u8>,
}

impl BitReader {
    fn new(buffer: VecDeque<u8>) -> Self {
        Self {
            bit: 0,
            offset: 0,
            buffer,
        }
    }

    fn read_bit(&mut self) -> (u8, bool) {
        if self.offset >= self.buffer.len() {
            return (0, false);
        }
        let b = self.buffer[self.offset];
        let k = (b >> (7 - self.bit)) & 1;
        self.bit = (self.bit + 1) % 8;
        if self.bit == 0 {
            self.offset += 1;
        }
        (k, true)
    }

    fn read_bits(&mut self, a: usize, b: &mut Vec<u8>) -> usize {
        b.clear();
        for _ in 0..a {
            let (k, ok) = self.read_bit();
            if !ok {
                return b.len();
            }
            b.push(k);
        }
        b.len()
    }
}

fn bits_to_byte(b: &[u8]) -> u8 {
    let mut k = 0;
    for i in b {
        k = (k << 1) + i;
    }
    return k;
}

fn bits_to_int(b: &[u8]) -> usize {
    let mut k = 0;
    for &i in b {
        k = (k << 1) + i as usize;
    }
    return k;
}

fn nibbles_to_int(b: &[u8]) -> usize {
    let mut k = 0;
    for &i in b {
        k = (k << 4) + i as usize;
    }
    return k;
}

fn parse_packet_header(tokens: &[usize]) -> Option<(usize, usize, &[usize])> {
    match tokens {
        [version, id, rest @ ..] => Some((*version, *id, rest)),
        _ => None,
    }
}

fn parse_packet_literal(tokens: &[usize]) -> Option<(usize, usize, &[usize])> {
    match tokens {
        [val, offset, rest @ ..] => Some((*val, *offset, rest)),
        _ => None,
    }
}

fn parse_subpacket_opts(tokens: &[usize]) -> Option<(usize, usize, usize, &[usize])> {
    match tokens {
        [mode, l, offset, rest @ ..] => Some((*mode, *l, *offset, rest)),
        _ => None,
    }
}

fn parse_subpackets_0(
    l: usize,
    offset: usize,
    tokens: &[usize],
) -> Option<(Vec<usize>, usize, &[usize])> {
    let mut vals = Vec::new();
    let mut cur_offset = offset;
    let mut cur_tokens = tokens;
    loop {
        if cur_offset - offset >= l {
            break;
        }
        if let Some((val, o, t)) = eval_packet(cur_tokens) {
            cur_offset = o;
            cur_tokens = t;
            vals.push(val);
        } else {
            return None;
        }
    }
    Some((vals, cur_offset, cur_tokens))
}

fn parse_subpackets_1(
    l: usize,
    offset: usize,
    tokens: &[usize],
) -> Option<(Vec<usize>, usize, &[usize])> {
    let mut vals = Vec::new();
    let mut cur_offset = offset;
    let mut cur_tokens = tokens;
    for _ in 0..l {
        if let Some((val, o, t)) = eval_packet(cur_tokens) {
            cur_offset = o;
            cur_tokens = t;
            vals.push(val);
        } else {
            return None;
        }
    }
    Some((vals, cur_offset, cur_tokens))
}

fn eval_packet(tokens: &[usize]) -> Option<(usize, usize, &[usize])> {
    let mut cur_offset;
    let mut cur_tokens = tokens;
    let id = match parse_packet_header(cur_tokens) {
        Some((_, id, t)) => {
            cur_tokens = t;
            id
        }
        None => return None,
    };
    if id == 4 {
        return parse_packet_literal(cur_tokens);
    }
    let (mode, l) = match parse_subpacket_opts(cur_tokens) {
        Some((mode, l, o, t)) => {
            cur_offset = o;
            cur_tokens = t;
            (mode, l)
        }
        None => return None,
    };
    let vals = match if mode == 0 {
        parse_subpackets_0(l, cur_offset, cur_tokens)
    } else {
        parse_subpackets_1(l, cur_offset, cur_tokens)
    } {
        Some((vals, o, t)) => {
            cur_offset = o;
            cur_tokens = t;
            vals
        }
        None => return None,
    };
    let k = match id {
        0 => vals.iter().sum::<usize>(),
        1 => vals.iter().product::<usize>(),
        2 => {
            if let Some(&v) = vals.iter().min() {
                v
            } else {
                return None;
            }
        }
        3 => {
            if let Some(&v) = vals.iter().max() {
                v
            } else {
                return None;
            }
        }
        5 => {
            if let [a, b] = vals[..] {
                if a > b {
                    1
                } else {
                    0
                }
            } else {
                return None;
            }
        }
        6 => {
            if let [a, b] = vals[..] {
                if a < b {
                    1
                } else {
                    0
                }
            } else {
                return None;
            }
        }
        7 => {
            if let [a, b] = vals[..] {
                if a == b {
                    1
                } else {
                    0
                }
            } else {
                return None;
            }
        }
        _ => return None,
    };
    Some((k, cur_offset, cur_tokens))
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let bitstream = if let Some(line) = reader.lines().flat_map(|i| i).next() {
        (0..line.len())
            .flat_map(|i| u8::from_str_radix(&line[i..i + 1], 16))
            .collect::<Vec<_>>()
            .chunks(2)
            .flat_map(|i| match i {
                &[a, b] => Some((a << 4) + b),
                &[a] => Some(a << 4),
                _ => None,
            })
            .collect::<VecDeque<u8>>()
    } else {
        return Err("Invalid line".into());
    };

    let mut tokens = Vec::new();
    let mut bits = BitReader::new(bitstream);
    let mut buf = Vec::with_capacity(15);
    let mut bit_offset = 0;
    let mut version_sum = 0;
    loop {
        if bits.read_bits(3, &mut buf) != 3 {
            break;
        }
        bit_offset += 3;
        let version = bits_to_byte(&buf[..]);
        if bits.read_bits(3, &mut buf) != 3 {
            break;
        }
        bit_offset += 3;
        let id = bits_to_byte(&buf[..]);
        tokens.push(version as usize);
        tokens.push(id as usize);
        version_sum += version as usize;
        if id == 4 {
            let mut nibbles = Vec::new();
            loop {
                if bits.read_bits(5, &mut buf) != 5 {
                    break;
                }
                bit_offset += 5;
                nibbles.push(bits_to_byte(&buf[1..]));
                if buf[0] == 0 {
                    break;
                }
            }
            tokens.push(nibbles_to_int(&nibbles));
        } else {
            if bits.read_bits(1, &mut buf) != 1 {
                break;
            }
            bit_offset += 1;
            let mode = buf[0];
            tokens.push(mode as usize);
            if mode == 0 {
                if bits.read_bits(15, &mut buf) != 15 {
                    break;
                }
                bit_offset += 15;
                tokens.push(bits_to_int(&buf[..]));
            } else {
                if bits.read_bits(11, &mut buf) != 11 {
                    break;
                }
                bit_offset += 11;
                tokens.push(bits_to_int(&buf[..]));
            }
        }
        tokens.push(bit_offset);
    }
    println!("Part 1: {}", version_sum);
    if let Some((val, _, _)) = eval_packet(&tokens[..]) {
        println!("Part 2: {}", val);
    } else {
        return Err("Failed eval".into());
    }
    Ok(())
}
