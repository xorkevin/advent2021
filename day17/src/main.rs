#[derive(Clone, Copy)]
struct Vec2(i32, i32);

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let (x1, x2, y1, y2) = (25, 67, -260, -200);
    let mut maxvy = 0;
    {
        let mut maxy = 0;
        let mut lower = 256;
        let mut upper = 265;
        while upper >= lower {
            let vy = lower + (upper - lower) / 2;
            if let Some(k) = simulate(Vec2(0, 0), Vec2(7, vy), x1, x2, y1, y2) {
                if k > maxy {
                    maxy = k;
                    maxvy = vy;
                }
                lower = vy + 1;
            } else {
                upper = vy - 1;
            }
        }
        println!("Part 1: {}", maxy)
    }
    {
        println!(
            "Part 2: {}",
            (7..x2 + 1)
                .into_iter()
                .flat_map(|i| (y1 - 1..maxvy + 1).into_iter().map(move |j| Vec2(i, j)))
                .filter_map(|p| simulate(Vec2(0, 0), p, x1, x2, y1, y2))
                .count()
        );
    }
    Ok(())
}

fn simulate(
    Vec2(mut x, mut y): Vec2,
    Vec2(mut vx, mut vy): Vec2,
    x1: i32,
    x2: i32,
    y1: i32,
    y2: i32,
) -> Option<i32> {
    let mut maxy = 0;
    loop {
        x += vx;
        y += vy;

        if vx > 0 {
            vx -= 1;
        } else if vx < 0 {
            vx += 1;
        }
        vy -= 1;

        if y > maxy {
            maxy = y;
        }

        if in_target(Vec2(x, y), x1, x2, y1, y2) {
            return Some(maxy);
        }
        if y < y1 {
            return None;
        }
        if vx == 0 && (x < x1 || x > x2) {
            return None;
        }
    }
}

fn in_target(Vec2(x, y): Vec2, x1: i32, x2: i32, y1: i32, y2: i32) -> bool {
    x >= x1 && x <= x2 && y >= y1 && y <= y2
}
