use std::collections::{HashMap, HashSet};
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

#[derive(PartialEq, Eq, Hash, Clone, Copy)]
struct Pos(usize, usize);

struct Graph {
    nodes: HashMap<String, bool>,
    edges: HashMap<String, Vec<String>>,
}

impl Graph {
    fn new() -> Self {
        Self {
            nodes: HashMap::new(),
            edges: HashMap::new(),
        }
    }

    fn is_upper(s: &str) -> bool {
        s == s.to_ascii_uppercase()
    }

    fn add_edge(&mut self, a: String, b: String) {
        self.nodes.insert(a.clone(), Graph::is_upper(&a));
        self.nodes.insert(b.clone(), Graph::is_upper(&b));
        self.edges
            .entry(a.clone())
            .or_insert(Vec::new())
            .push(b.clone());
        self.edges.entry(b).or_insert(Vec::new()).push(a);
    }

    fn find_path_rec(&self, start: &str, end: &str, path: &mut Vec<String>) -> usize {
        if start == end {
            return 1;
        }
        let visited = path.iter().map(|i| i.clone()).collect::<HashSet<_>>();
        self.edges[start]
            .iter()
            .filter_map(|i| {
                if visited.contains(i) && !self.nodes[i] {
                    None
                } else {
                    path.push(i.clone());
                    let k = self.find_path_rec(&i, end, path);
                    path.pop();
                    Some(k)
                }
            })
            .sum()
    }

    fn find_path(&self, start: &str, end: &str) -> usize {
        self.find_path_rec(start, end, &mut vec![start.into()])
    }

    fn find_path_rec2(&self, start: &str, end: &str, path: &mut Vec<String>) -> usize {
        if start == end {
            return 1;
        }
        let (has_twice, visited) = path.iter().fold((false, HashSet::new()), |(b, mut s), i| {
            let n = b || s.contains(i) && !self.nodes[i];
            s.insert(i.clone());
            (n, s)
        });
        self.edges[start]
            .iter()
            .filter_map(|i| {
                if visited.contains(i) && !self.nodes[i] {
                    if has_twice || i == "start" || i == "end" {
                        return None;
                    }
                }
                path.push(i.clone());
                let k = self.find_path_rec2(&i, end, path);
                path.pop();
                Some(k)
            })
            .sum()
    }

    fn find_path2(&self, start: &str, end: &str) -> usize {
        self.find_path_rec2(start, end, &mut vec![start.into()])
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let graph = reader
        .lines()
        .flat_map(|i| i)
        .filter_map(|i| {
            if let [a, b] = i.splitn(2, "-").collect::<Vec<_>>()[..] {
                Some((a.to_owned(), b.to_owned()))
            } else {
                None
            }
        })
        .fold(Graph::new(), |mut g, (a, b)| {
            g.add_edge(a, b);
            g
        });

    println!("Part 1: {}", graph.find_path("start", "end"));
    println!("Part 2: {}", graph.find_path2("start", "end"));
    Ok(())
}
