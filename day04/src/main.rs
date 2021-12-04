use std::collections::hash_map::Entry;
use std::collections::HashMap;
use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

struct Pos(usize, usize);

struct Board {
    width: usize,
    height: usize,
    unmarked: HashMap<i32, Pos>,
    marked: HashMap<i32, Pos>,
    rows: Vec<usize>,
    cols: Vec<usize>,
}

impl Board {
    fn new() -> Self {
        Board {
            width: 0,
            height: 0,
            unmarked: HashMap::new(),
            marked: HashMap::new(),
            rows: Vec::new(),
            cols: Vec::new(),
        }
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);
    let lines = reader.lines();

    let mut nums = Vec::new();
    let mut boards = Vec::new();
    let mut board = Board::new();
    let mut first = true;
    for line in lines {
        if first {
            first = false;
            nums = line?
                .split(",")
                .flat_map(|num| num.parse::<i32>())
                .collect::<Vec<_>>();
            continue;
        }
        let line = line?;
        if line == "" {
            if board.height > 0 {
                board.rows = vec![0; board.height];
                board.cols = vec![0; board.width];
                boards.push(board);
                board = Board::new();
            }
            continue;
        }
        let h = board.height;
        board.height += 1;
        let mut w = 0;
        for (n, i) in line.split_ascii_whitespace().enumerate() {
            w += 1;
            if let Entry::Vacant(v) = board.unmarked.entry(i.parse::<i32>()?) {
                v.insert(Pos(n, h));
            } else {
                return Err("Duplicate num".into());
            }
        }
        board.width = w;
    }
    if board.height > 0 {
        board.rows = vec![0; board.height];
        board.cols = vec![0; board.width];
        boards.push(board);
    }

    let mut skip_map = HashSet::new();
    first = true;
    for i in nums {
        let l = boards.len();
        for (n, j) in boards.iter_mut().enumerate() {
            if skip_map.contains(&n) {
                continue;
            }
            let k = mark_board(j, i);
            if k > -1 {
                if first {
                    first = false;
                    println!("Part 1: {}", k * i);
                }
                skip_map.insert(n);
                if skip_map.len() >= l {
                    println!("Part 2: {}", k * i);
                    return Ok(());
                }
            }
        }
    }

    Ok(())
}

fn mark_board(board: &mut Board, num: i32) -> i32 {
    if let Entry::Occupied(e) = board.unmarked.entry(num) {
        let (k, v) = e.remove_entry();
        let Pos(x, y) = v;
        board.marked.insert(k, v);
        board.rows[y] += 1;
        board.cols[x] += 1;
        if board.rows[y] >= board.height {
            return score_board(board);
        }
        if board.cols[x] >= board.width {
            return score_board(board);
        }
    }
    -1
}

fn score_board(board: &Board) -> i32 {
    board.unmarked.keys().sum()
}
