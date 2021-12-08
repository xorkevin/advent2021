use std::collections::{HashMap, HashSet};
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

struct Constraints {
    assign: HashMap<u8, i32>,
    opts: HashMap<u8, HashSet<i32>>,
    common6: HashSet<u8>,
    common5: HashSet<u8>,
}

fn full_opts() -> HashSet<i32> {
    HashSet::from_iter(0..7)
}

impl Constraints {
    fn new() -> Self {
        Constraints {
            assign: HashMap::new(),
            opts: HashMap::from_iter((0..7).into_iter().map(|i| (b'a' + i, full_opts()))),
            common6: HashSet::from_iter((0..7).into_iter().map(|i| (b'a' + i))),
            common5: HashSet::from_iter((0..7).into_iter().map(|i| (b'a' + i))),
        }
    }

    fn reduce_opts(&mut self, wires: &[u8], segs: &[i32]) {
        let mut changed = false;
        for i in wires {
            for j in segs {
                if let Some(v) = self.opts.get_mut(i) {
                    if v.contains(j) {
                        changed = true;
                        v.remove(j);
                    }
                }
            }
        }
        if !changed {
            return;
        }
        loop {
            let mut changed = false;
            for (&k, v) in self.opts.iter_mut() {
                if self.assign.contains_key(&k) {
                    continue;
                }
                if v.len() == 1 {
                    changed = true;
                    if let Some(&i) = v.iter().next() {
                        self.assign.insert(k, i);
                    }
                    continue;
                }
                for i in self.assign.values() {
                    if v.contains(i) {
                        changed = true;
                        v.remove(i);
                    }
                }
            }
            if !changed {
                return;
            }
        }
    }

    fn reduce_common6(&mut self, wires: &[u8]) {
        let s = HashSet::from_iter(wires.iter().map(|&i| i));
        self.common6 = self.common6.intersection(&s).map(|&i| i).collect();
    }

    fn reduce_common5(&mut self, wires: &[u8]) {
        let s = HashSet::from_iter(wires.iter().map(|&i| i));
        self.common5 = self.common5.intersection(&s).map(|&i| i).collect();
    }

    fn translate_wires(&self, wires: &[u8]) -> Vec<i32> {
        wires.iter().map(|i| self.assign[i]).collect()
    }

    fn translate_seg6(&self, wires: &[u8]) -> i32 {
        let mut missing2 = true;
        let mut missing3 = true;
        let mut missing4 = true;
        for i in self.translate_wires(wires) {
            match i {
                2 => missing2 = false,
                3 => missing3 = false,
                4 => missing4 = false,
                _ => (),
            }
        }
        if missing3 {
            return 0;
        }
        if missing2 {
            return 6;
        }
        if missing4 {
            return 9;
        }
        panic!("Should not reach");
    }

    fn translate_seg5(&self, wires: &[u8]) -> i32 {
        let mut has1 = false;
        let mut has2 = false;
        let mut has4 = false;
        let mut has5 = false;
        for i in self.translate_wires(wires) {
            match i {
                1 => has1 = true,
                2 => has2 = true,
                4 => has4 = true,
                5 => has5 = true,
                _ => (),
            }
        }
        if has2 && has4 {
            return 2;
        }
        if has2 && has5 {
            return 3;
        }
        if has1 && has5 {
            return 5;
        }
        panic!("Should not reach");
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let segs1_inv = vec![0, 1, 3, 4, 6];
    let segs4_inv = vec![0, 4, 6];
    let segs7_inv = vec![1, 3, 4, 6];
    let segs6_inv = vec![2, 3, 4];
    let segs5_inv = vec![1, 2, 4, 5];

    let mut count = 0;
    let mut count2 = 0;
    for line in reader.lines() {
        if let [cases, output] = line?.splitn(2, " | ").collect::<Vec<_>>()[..] {
            let mut constraints = Constraints::new();
            for i in cases.split_ascii_whitespace() {
                let wires = i.as_bytes();
                let l = wires.len();
                match l {
                    2 => constraints.reduce_opts(wires, &segs1_inv),
                    4 => constraints.reduce_opts(wires, &segs4_inv),
                    3 => constraints.reduce_opts(wires, &segs7_inv),
                    6 => constraints.reduce_common6(wires),
                    5 => constraints.reduce_common5(wires),
                    _ => (),
                }
            }
            if constraints.common6.len() == 4 {
                constraints.reduce_opts(
                    &constraints.common6.iter().map(|&i| i).collect::<Vec<_>>()[..],
                    &segs6_inv,
                );
            }
            if constraints.common5.len() == 3 {
                constraints.reduce_opts(
                    &constraints.common5.iter().map(|&i| i).collect::<Vec<_>>()[..],
                    &segs5_inv,
                );
            }
            if constraints.assign.len() != 7 {
                return Err("Failed to assign all".into());
            }
            let (a1, a2) = output.split_ascii_whitespace().fold((0, 0), |(a1, a2), i| {
                let wires = i.as_bytes();
                let l = wires.len();
                match l {
                    2 => (a1 + 1, a2 * 10 + 1),
                    4 => (a1 + 1, a2 * 10 + 4),
                    3 => (a1 + 1, a2 * 10 + 7),
                    7 => (a1 + 1, a2 * 10 + 8),
                    6 => (a1, a2 * 10 + constraints.translate_seg6(wires)),
                    5 => (a1, a2 * 10 + constraints.translate_seg5(wires)),
                    _ => (a1, a2),
                }
            });
            count += a1;
            count2 += a2;
        }
    }
    println!("Part 1: {}", count);
    println!("Part 2: {}", count2);
    Ok(())
}
