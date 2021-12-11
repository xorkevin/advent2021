use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

#[derive(PartialEq, Eq, Hash, Clone, Copy)]
struct Pos(usize, usize);

struct Grid {
    rows: Vec<Vec<i32>>,
    w: usize,
    h: usize,
}

impl Grid {
    fn new(rows: Vec<Vec<i32>>) -> Self {
        let h = rows.len();
        let w = if h > 0 { rows[0].len() } else { 0 };
        Self { rows, w, h }
    }

    fn step(&mut self) -> usize {
        let mut flashes = HashSet::new();
        let mut queue = Vec::new();
        for i in 0..self.h {
            for j in 0..self.w {
                self.rows[i][j] += 1;
                if self.rows[i][j] > 9 {
                    let k = Pos(j, i);
                    flashes.insert(k);
                    queue.push(k);
                }
            }
        }
        loop {
            let Pos(x, y) = match queue.pop() {
                Some(i) => i,
                None => break,
            };
            if self.incr(x.wrapping_sub(1), y, &flashes) {
                let k = Pos(x.wrapping_sub(1), y);
                flashes.insert(k);
                queue.push(k);
            }
            if self.incr(x.wrapping_sub(1), y.wrapping_sub(1), &flashes) {
                let k = Pos(x.wrapping_sub(1), y.wrapping_sub(1));
                flashes.insert(k);
                queue.push(k);
            }
            if self.incr(x, y.wrapping_sub(1), &flashes) {
                let k = Pos(x, y.wrapping_sub(1));
                flashes.insert(k);
                queue.push(k);
            }
            if self.incr(x + 1, y.wrapping_sub(1), &flashes) {
                let k = Pos(x + 1, y.wrapping_sub(1));
                flashes.insert(k);
                queue.push(k);
            }
            if self.incr(x + 1, y, &flashes) {
                let k = Pos(x + 1, y);
                flashes.insert(k);
                queue.push(k);
            }
            if self.incr(x + 1, y + 1, &flashes) {
                let k = Pos(x + 1, y + 1);
                flashes.insert(k);
                queue.push(k);
            }
            if self.incr(x, y + 1, &flashes) {
                let k = Pos(x, y + 1);
                flashes.insert(k);
                queue.push(k);
            }
            if self.incr(x.wrapping_sub(1), y + 1, &flashes) {
                let k = Pos(x.wrapping_sub(1), y + 1);
                flashes.insert(k);
                queue.push(k);
            }
        }
        for &Pos(x, y) in &flashes {
            self.rows[y][x] = 0;
        }
        flashes.len()
    }

    fn incr(&mut self, x: usize, y: usize, flashes: &HashSet<Pos>) -> bool {
        if self.out_bounds(x, y) {
            return false;
        }
        self.rows[y][x] += 1;
        if self.rows[y][x] > 9 {
            !flashes.contains(&Pos(x, y))
        } else {
            false
        }
    }

    fn out_bounds(&self, x: usize, y: usize) -> bool {
        x >= self.w || y >= self.h
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut grid = Grid::new(
        reader
            .lines()
            .flat_map(|i| i)
            .map(|i| {
                i.as_bytes()
                    .into_iter()
                    .map(|&j| j as i32 - b'0' as i32)
                    .collect()
            })
            .collect(),
    );
    println!(
        "Part 1: {}",
        (0..100).into_iter().map(|_| grid.step()).sum::<usize>()
    );
    let total = grid.w * grid.h;
    for i in 101.. {
        if grid.step() == total {
            println!("Part 2: {}", i);
            break;
        }
    }
    Ok(())
}
