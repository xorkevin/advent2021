use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let nums = {
        let mut nums = Vec::new();
        for line in reader.lines() {
            nums.push(line?.into_bytes());
        }
        nums
    };

    if nums.len() == 0 {
        return Ok(());
    }

    let numbits = nums[0].len();

    let mut count = vec![0; numbits];
    for i in &nums {
        for (n, &j) in i.iter().enumerate() {
            if j == b'1' {
                count[n] += 1;
            } else {
                count[n] -= 1;
            }
        }
    }
    let mut max = vec![0; numbits];
    let mut min = vec![0; numbits];
    for (n, &i) in count.iter().enumerate() {
        if i > 0 {
            max[n] = b'1';
            min[n] = b'0';
        } else {
            max[n] = b'0';
            min[n] = b'1';
        }
    }
    println!("Part 1: {}", btoi(&max[..]) * btoi(&min[..]));

    let mut most = nums.iter().collect::<Vec<_>>();
    let mut least = nums.iter().collect::<Vec<_>>();

    for i in 0..numbits {
        if most.len() < 2 {
            break;
        }
        most = find_common(true, i, most);
    }
    for i in 0..numbits {
        if least.len() < 2 {
            break;
        }
        least = find_common(false, i, least);
    }

    if most.len() == 1 && least.len() == 1 {
        println!("Part 2: {}", btoi(&most[0][..]) * btoi(&least[0][..]));
    }
    Ok(())
}

fn find_common<'a>(most: bool, pos: usize, nums: Vec<&'a Vec<u8>>) -> Vec<&'a Vec<u8>> {
    let mut onecounts = 0;
    let mut zerocounts = 0;
    for i in &nums {
        if i[pos] == b'1' {
            onecounts += 1;
        } else {
            zerocounts += 1;
        }
    }

    let zeros = if most {
        zerocounts > onecounts
    } else {
        zerocounts <= onecounts
    };

    let mut res = Vec::new();

    for i in nums {
        if zeros {
            if i[pos] == b'0' {
                res.push(i);
            }
        } else {
            if i[pos] == b'1' {
                res.push(i);
            }
        }
    }

    res
}

fn btoi(b: &[u8]) -> i32 {
    let mut c = 0;
    for &i in b {
        c = (c << 1) + (if i == b'1' { 1 } else { 0 });
    }
    return c;
}
