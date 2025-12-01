import math

def read_input(fpath):
    lines = None
    with open(fpath) as inf:
        lines = inf.readlines()
    return lines


def parse_rotation(rotation):
    sign = 1.0
    if rotation[0] == "L":
        sign = -1.0
   
    return sign * int(rotation[1:])

def rotate(dial, rotation):
    return (dial + rotation) % 100

def sign(a):
    if (a < 0):
        return -1
    elif (a > 0):
        return 1
    else:
        return 0
   

def count_turns_bruteforce(dial, rotation):
    s = sign(rotation)
    rot = int(math.fabs(rotation))
    count = 0

    for r in range(rot):
        dial = rotate(dial, s)
        if dial == 0:
            count += 1    
    return count


def part1(data):
    dial = 50
    count = 0
    for line in data:
        rotation = parse_rotation(line)
        dial = rotate(dial, rotation)
        if (dial == 0):
            count += 1
    return count

def part2(data):
    dial = 50
    count = 0
    for line in data:
        rotation = parse_rotation(line)
        turns_b = count_turns_bruteforce(dial, rotation)                   
        dial = rotate(dial, rotation)
        count += turns_b

    return count


data = read_input("input1.txt")
print("Part 1:", part1(data))
print("Part 2:", part2(data))