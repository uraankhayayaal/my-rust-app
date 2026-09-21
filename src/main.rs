use rand::Rng;
use std::cmp::Ordering;
use std::io;

extern crate rand;

fn max_value_from(difficulty: &str) -> u32 {
    match difficulty {
        "easy" => 10,
        "hard" => 1000,
        _ => 100,
    }
}

fn main() {
    println!("Guess the number!");

    let difficulty = std::env::args().nth(1).unwrap_or_else(|| "normal".to_string());
    let max_value = max_value_from(&difficulty);
    let secret_number = rand::rng().random_range(1..=max_value);

    println!("The secret number is: {secret_number} (difficulty: {difficulty}, max: {max_value})");

    loop {
        println!("Please input your guess.");

        let mut guess = String::new();

        io::stdin()
            .read_line(&mut guess)
            .expect("Failed to read line");

        // let guess: u32 = guess.trim().parse().expect("Please type a number!");
        let guess: u32 = match guess.trim().parse() {
            Ok(num) => num,
            Err(_) => continue,
        };

        println!("You guessed: {guess}");

        match guess.cmp(&secret_number) {
            Ordering::Less => println!("Too small!"),
            Ordering::Greater => println!("Too big!"),
            Ordering::Equal => {
                println!("You win!");
                break;
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::max_value_from;

    #[test]
    fn easy_difficulty_bounds_to_ten() {
        assert_eq!(max_value_from("easy"), 10);
    }

    #[test]
    fn hard_difficulty_bounds_to_thousand() {
        assert_eq!(max_value_from("hard"), 1000);
    }

    #[test]
    fn unknown_difficulty_falls_back_to_hundred() {
        assert_eq!(max_value_from("weird"), 100);
        assert_eq!(max_value_from(""), 100);
    }
}