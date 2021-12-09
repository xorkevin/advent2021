use std::cmp::Reverse;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

struct Grid {
    rows: Vec<Vec<u8>>,
    h: usize,
    w: usize,
    basin: Vec<Vec<bool>>,
}

impl Grid {
    fn new(rows: Vec<Vec<u8>>) -> Self {
        let h = rows.len();
        let w = rows[0].len();
        Self {
            rows,
            h,
            w,
            basin: (0..h).into_iter().map(|_| vec![false; w]).collect(),
        }
    }

    fn out_bounds(&self, x: usize, y: usize) -> bool {
        x >= self.w || y >= self.h
    }

    fn in_bounds(&self, x: usize, y: usize) -> bool {
        !self.out_bounds(x, y)
    }

    fn is_low(&self, x: usize, y: usize) -> bool {
        let v = self.rows[y][x];
        if self.in_bounds(x.wrapping_sub(1), y) && self.rows[y][x.wrapping_sub(1)] <= v {
            return false;
        }
        if self.in_bounds(x, y.wrapping_sub(1)) && self.rows[y.wrapping_sub(1)][x] <= v {
            return false;
        }
        if self.in_bounds(x + 1, y) && self.rows[y][x + 1] <= v {
            return false;
        }
        if self.in_bounds(x, y + 1) && self.rows[y + 1][x] <= v {
            return false;
        }
        true
    }

    fn mark_basin(&mut self, x: usize, y: usize) -> i32 {
        if !self.in_bounds(x, y) {
            return 0;
        }
        if self.rows[y][x] == b'9' {
            return 0;
        }
        if self.basin[y][x] {
            return 0;
        }
        self.basin[y][x] = true;
        1 + self.mark_basin(x.wrapping_sub(1), y)
            + self.mark_basin(x, y.wrapping_sub(1))
            + self.mark_basin(x + 1, y)
            + self.mark_basin(x, y + 1)
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut grid = Grid::new(
        reader
            .lines()
            .flat_map(|i| i)
            .map(|i| Vec::from(i.as_bytes()))
            .collect(),
    );
    let w = grid.w;
    let mut sizes = (0..grid.h)
        .into_iter()
        .flat_map(|r| (0..w).into_iter().map(move |c| (r, c)))
        .flat_map(|(r, c)| {
            if grid.is_low(c, r) {
                Some((
                    (grid.rows[r][c] as i32) - (b'0' as i32) + 1,
                    grid.mark_basin(c, r),
                ))
            } else {
                None
            }
        })
        .collect::<Vec<_>>();
    sizes.sort_by_key(|&(_, a)| Reverse(a));
    println!("Part 1: {}", sizes.iter().map(|&(a, _)| a).sum::<i32>());
    println!(
        "Part 2: {}",
        sizes.iter().take(3).map(|&(_, a)| a).product::<i32>()
    );
    Ok(())
}
