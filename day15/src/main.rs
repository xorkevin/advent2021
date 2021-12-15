use std::cmp::Ordering;
use std::collections::{BinaryHeap, HashSet};
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

#[derive(Clone, Copy, Hash, Eq, PartialEq, PartialOrd, Ord, Debug)]
struct Point(usize, usize);
#[derive(Eq, PartialEq, Debug)]
struct Item(Point, usize, usize);

impl PartialOrd for Item {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl Ord for Item {
    fn cmp(&self, Item(v, g, f): &Self) -> Ordering {
        let Item(sv, sg, sf) = self;
        f.cmp(sf).then_with(|| g.cmp(sg)).then_with(|| sv.cmp(v))
    }
}

fn abs(a: usize, b: usize) -> usize {
    if a < b {
        b - a
    } else {
        a - b
    }
}

fn manhattan(Point(ax, ay): Point, Point(bx, by): Point) -> usize {
    abs(ax, bx) + abs(ay, by)
}

fn neighbors(Point(x, y): Point) -> Vec<Point> {
    vec![
        Point(x.wrapping_sub(1), y),
        Point(x, y.wrapping_sub(1)),
        Point(x + 1, y),
        Point(x, y + 1),
    ]
}

fn in_bounds(Point(x, y): Point, w: usize, h: usize) -> bool {
    x < w && y < h
}

fn getval(grid: &Vec<Vec<usize>>, Point(x, y): Point, w: usize, h: usize) -> usize {
    (grid[y % h][x % w] + x / w + y / h - 1) % 9 + 1
}

fn pathfind(
    grid: &Vec<Vec<usize>>,
    vw: usize,
    vh: usize,
    start: Point,
    end: Point,
) -> Option<usize> {
    if grid.len() == 0 {
        return None;
    }
    let w = grid[0].len();
    let h = grid.len();
    let mut openset = BinaryHeap::new();
    openset.push(Item(start, 0, manhattan(start, end)));
    let mut closedset = HashSet::new();
    while let Some(Item(v, g, _)) = openset.pop() {
        if closedset.contains(&v) {
            continue;
        }
        closedset.insert(v);
        if v == end {
            return Some(g);
        }
        for i in neighbors(v) {
            if !in_bounds(i, vw, vh) || closedset.contains(&i) {
                continue;
            }
            let ng = g + getval(grid, i, w, h);
            openset.push(Item(i, ng, ng + manhattan(i, end)));
        }
    }
    None
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let grid = reader
        .lines()
        .flat_map(|i| i)
        .map(|line| {
            line.as_bytes()
                .into_iter()
                .map(|&i| (i - b'0') as usize)
                .collect::<Vec<_>>()
        })
        .collect::<Vec<_>>();
    println!(
        "Part 1: {}",
        pathfind(
            &grid,
            grid[0].len(),
            grid.len(),
            Point(0, 0),
            Point(grid[0].len() - 1, grid.len() - 1)
        )
        .ok_or("Fail part 1")?
    );
    println!(
        "Part 2: {}",
        pathfind(
            &grid,
            grid[0].len() * 5,
            grid.len() * 5,
            Point(0, 0),
            Point(grid[0].len() * 5 - 1, grid.len() * 5 - 1)
        )
        .ok_or("Fail part 2")?
    );
    Ok(())
}
