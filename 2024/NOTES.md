# Day 1

Day 1 was fairly straight forward. It was a matter of parsing the 2 lists that 
were given as vertical columns instead of 2 lines. Then it was a matter of sorting 
those 2 lists to get the distance between each pair in the 2 lists. I decided to write
quicksort by hand since I had it memorized and to be honest, I haven't used Go's generic
library functions since I had stopped using Go before generics landed. 

Part 2 of day 1 was similarly straightforward since it was just a matter of counting occurrances
of each element in list 1 with list 2 using a HashMap.

# Day 2

The problem was a list of list of numbers. Each number in the list is a level and a list is a report of levels. 
The report can be considered safe or not depending on how much each level increases or decreases within the report.
Finding the unsafe levels in part 1 was fairly straight forward. 

In part 2, you could fix unsafe reports by removing a level from the list. I tried to be smart at first by finding the
I first tried to figure out the direction and iterate on the report and if I found a bad level tried to remove it. But
that solution didn't work because some reports could be fixed by removing the first or second level that would change
the direction of the report by doing so. There were several edge cases like that, so I ended up giving up on the smarter
solution and went with the brute force solution of trying to fix each unsafe report by removing a level at a time and 
testing the report. If it succeeded I went with it. The solution turned out to run pretty well with the input given, 
so I went with it. 

```go
func IsSafeWithDampenerBruteForce(report []int) bool {
	if IsSafe(report) {
		return true
	} else {
		for i := range report {
			var dampenedReport []int = make([]int, len(report))
			copy(dampenedReport, report)
			dampenedReport = append(dampenedReport[:i], dampenedReport[i+1:]...)

			if IsSafe(dampenedReport) {
				return true
			}
		}
		return false
	}
}
```

Lesson of the day. Sometimes try the brute force method first, it might surprise you at how effective it is.

# Day 3

Day 3 was maybe one of my favorites since it was a parsing challenge. Given a bunch of garbled text containing some
correct instructions `mul(x,x)`, parse the correct instruction and run them. 

e.g. input `xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))`

Part 2 simply added 2 new instructions `do()` and `don't()` which toggle whether you should execute the instructions 
following. It didn't really add to the challenge.

For this I decided to try out the more interesting solution of writing a custom lexer/tokenizer and parser for 
the instructions. Then I compared it to a solution using regex.

The custom lexer and parser followed the pattern described by Rob Pike in his talk [Lexical Scanning in Go](https://www.youtube.com/watch?v=HxaD_trXwRE) ([Slides](https://go.dev/talks/2011/lex.slide#1)).
In it he describes a pattern of having the state machine that drives the lexer be a bunch of functions that return 
a function with the next state to be scanned. It looks something like this:

```go
type Lexer struct {
	input string
	start int
	pos   int
	items []Item
}

type stateFn func(l *Lexer) stateFn


func Lex(input string) []Item {
	lexer := &Lexer{
		input: input,
	}

	for state := lexCorrupted(lexer); ; {
		if lexer.pos >= len(input) {
			break
		}
		state = state(lexer)
	}
	return lexer.items
}

func lexCorrupted(l *Lexer) stateFn {
	if l.input[l.pos] == 'm' {
		return lexMul
	} else if l.input[l.pos] == 'd' {
		return lexDoDont
	} else {
		l.pos += 1
		l.start = l.pos
		return lexCorrupted
	}
}

func lexMul(l *Lexer) stateFn {
	if l.pos+4 <= len(l.input) && l.input[l.pos:l.pos+4] == "mul(" {
		l.items = append(l.items, Item{typ: itemMulStart, val: string(l.input[l.pos : l.pos+4])})
		l.pos += 4
		l.start = l.pos
		return lexNumber
	} else {
		l.pos += 1
		l.start = l.pos
		return lexCorrupted
	}
}
```

In his he returns the tokens with a channel, I decided to just keep the state in the lexer struct. I would have to think
more on what improvements returning the tokens through a channel and running the lexer in a goroutine would bring, but 
I found keeping the tokens worked for me here.

It was a fun process writing the lexer and the solution ended up pretty simple and correct on the first time. That makes
me smile when that happens. 

The regex version turned out to be much more succinct. Not a surprise there. My first stab at it used multiple regex
patterns which I then take the matches and merge them back together and sort on the index of the matched text. 
But I eventually found a better way using multiple patterns and the alternation operator.I don't often use more complex regexes,
so it was good to be reminded of the alternation operator to match multiple patterns and named groups to extract 
which pattern a given match had matched. 

Finally I benchmarked the various solutions which came out as follows:
```bash
> go test -bench=. -benchtime=5s ./day3/...
goos: darwin                                                            
goarch: arm64
pkg: jonfk.ca/advent-of-code/2024/day3
cpu: Apple M4
BenchmarkParse-10                          89307             66544 ns/op
BenchmarkParseWithNamedRegex-10             6637            912917 ns/op
BenchmarkParseWithMultiRegex-10            19225            314754 ns/op
PASS
ok      jonfk.ca/advent-of-code/2024/day3       22.133s
```

The custom lexer and parser was much faster, atcually by at least 1 order of magnitude.

Lessons learned:
- Regexes can be much more concise solutions but also much less performant vs a custom parser.
- Writing a custom lexer and parser can be much simpler than you might expect.

# Day 4

Day 4 was a grid of text inside of which you needed to find the word XMAS. I could have parsed the text into a matrix
of characters and at every point I found the start of XMAS, do a check in every direction forward, backwards, and diagonally
in all directions. But I found that the edge cases for that would be kind of boring to write and deal with and test, so
my solution for part 1 was to collect all the different directions strings could be in and do a check on the string 
whether `XMAS` or `SAMX` existed inside and count them. That turned out to be a little more difficult than I expected 
because I didn't quite understand all the directions the strings could go in. It turned out I needed to do a matrix
rotation to collect all diagonal strings which took some whiteboarding. 

Part 2 turned out to be much simpler since the pattern was much stricter using 2 MAS in an X pattern in any direction.
For that one I actually used a giant if statement with an iteration on the matrix I originally avoided in part 1.

Here is the monstruosity I wrote for part 2:

```go
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] == 'A' {
				if ((i-1 >= 0 && j+1 < len(matrix[i]) && matrix[i-1][j+1] == 'S' && i+1 < len(matrix) && j-1 >= 0 && matrix[i+1][j-1] == 'M') || (i-1 >= 0 && j+1 < len(matrix[i]) && matrix[i-1][j+1] == 'M' && i+1 < len(matrix) && j-1 >= 0 && matrix[i+1][j-1] == 'S')) && ((i-1 >= 0 && j-1 >= 0 && matrix[i-1][j-1] == 'S' && i+1 < len(matrix) && j+1 < len(matrix[i]) && matrix[i+1][j+1] == 'M') || (i-1 >= 0 && j-1 >= 0 && matrix[i-1][j-1] == 'M' && i+1 < len(matrix) && j+1 < len(matrix[i]) && matrix[i+1][j+1] == 'S')) {
					count += 1
				}
			}
		}
	}
```

# Day 5

Day 5's problem was a set of rules for which number could come before which number and lists of numbers. Part 1 was 
fairly straight forward since it only required checking whether a list followed the rules or not. To check that,
I inverted the rules by storing all the numbers that should come before each particular number in a `Map[int][]int`.
Then I could simply check if any number after a particular number existing in it's entry, that list would be incorrect.

Part 2 was much trickier since it required fixing the incorrect lists in the right order. I initially thought to sort
them using the rules. That turned out not to work out very well. I think because I implemented it using quicksort which
uses a comparison function that used the rules to check whether a particular number should be less than or equal to 
the pivot and the rules actually included cycles which the problem input carefully avoided the cycles within the lists.
I ended up needing to check reddit, where I saw the meme about non-transitivity of the comparison and that put me on
the right path. I think if I had used bubble sort like most people and simply swap the numbers that violate a rule, I 
would not have encountered that problem. I am still unsure whether that was the most optimal solution.

# Day 6

The problem was a map that contains empty spots and spots with obstacles, a guard that has a direction and changes direction
when it encounters an obstacle. I encoded the problem as a matrix and the cells and directions as enums. 

The first part was to find all distinct positions the guard would be at given a starting direction and position. 

The second part was to add an obstacle that would make the guard go in a loop. I initially tried to find a smart way
to detect whether the guard would be in a loop without actually having to compute the whole path the guard would take,
but that didn't seem to work.

After spending a few hours on my broken solution, I cleaned up what I had and started part 2 frorm scratch.I ended up 
implementing the simplest most bruteforce-ish solution I could think of. Part 1 already gave me a way to walk the guard
around, I simply added a loop detection statement to it that would break out of the loop and return if there is a loop,
then tried to add an obstacle to every position the guard had taken up on his original walk and try the run from it's 
initial position using this new matrix. That ended up work. 

I am still unsure why my optimal solution didn't work. But it speaks to implementing the simplest solution that could 
give you an answer first before attempting to make it more efficient. Lesson learned.

I probably should revisit this one and maybe try to make that more efficient solution work. Here is some pseudo code
of how it would work.

```
var visitedCells [][]bool // a given coordinate is true if the guard has been there on his original walk
var visitedDirections [][][]direction // The direction the guard has taken at this coordinate

for each coordinate x,y in visitedCells:
	simulate an obstacle at the next coordinate the guard would step using the direction he had taken in visitedDirections
	simulate walk, if guard steps back on the x,y with the direction taken, loop has been detected
	
```

# Day 7

Day 7 was another parsing and calculator problem. The solution to part 1 was fairly straight forward of parsing some 
numbers and a target number to get to with some operators. I used an enum to model the operators (+ and *) and use 
a permutation generation algorithm that I have down pretty well by now. To do the evaluation itself, I used a stack 
as I would to create an RPN (reverse polish notation) calculator.

The second part was simply adding a concat operator which I was lazy and implemented using string concatenation instead
of mathematically. I should think about how to do it mathematically and benchmark if there is any difference in performance.

Leaving a TODO here for that. I am really curious since I think the string concatenation version is actually pretty
decently efficient?

# Day 8

Day 8 was essentially a sort of graphing problem where you have points of different types, each pair of point of a 
particular type generate points. The first part the points generated from each pair must be distance equal to the 
distance between the points from the points. The second part, the points were on every point on the invisible line 
between the points. In truth, it was not the most difficult, but I took the longest here because it took me forever to
actually understand what the problem actually was and what to implement. 

# Day 9

The problem today was a representation of blocks in a file system which can be empty or contain blocks from a file. The
problem was to compact the blocks and calculate a checksum on the resulting blocks. Part 1 was fairly straight forward
as usual. 

But part 2 strained my reading comprehension. I misunderstood the order of the algorithm and ended up implementing 
a wrong algorithm on first try. On second try, I did end up with much nicer code. I ended up implementing some unnecessary
work because of my wrong first pass. I ended up debugging for a while not getting the solution but found out in the end
that my checksumming function had an assumption from part 1 that I didn't revisit which was causing the bug.

# Day 10

The problem for day 10 was a 2d map encoded as numbers that allows traversal from 0 to 9 in increments of 1. Part 1 
was actually simple enough to solve with a recursive solution directly with the input. Part 2 was similar but instead
of simply finding out how many ends I can find per start, I needed to find how many distinct paths. So I had to implement
the recursive solution in the iterative version. It was pretty fun actually because my first version for part 1 and 3rd 
version for part 2 worked. I had very little debugging to do. 

I also decided to write in Java for this one since I am interviewing recently and would like the practice in case some
companies are really strict about the programming language they want to use in the interview. This turned out to be 
the perfect problem to do that with since I wasted a lot of time on the tooling and setup. I decided to go with Jbang 
to run the Java solution because of it's simplicity. I just really did not want to debug a gradle build or even a maven
xml soup. I would have to admit that the build tooling is probably 40% of the reason why I don't do Java side projects
most of the time. The other is how verbose everything is. But I got to say that since Java 9, which was what I was using 
the last company I used Java on the job, Java has improved a lot. A few features I like a lot:

- record classes
- var for local variable type inference
- switch expressions
- pattern matching

# Day 11

The problem was deceptively simple. Part 1 was easy to do with a bruteforce solution. Part 2 was just part 1 with a 
much higher iteration number which wouldn't work with the brute force solution. I didn't get much time to think about 
the problem since I have been taking care of Alexander these days, So I checked out the subreddit where I saw the 
suggestion of using a recursive solution with memoization. Took me a bit of time to think about it but the solution 
ended up pretty simple too. I still had issues with number overflowing which required me to use long instead of int 
for the result.

# Day 12

Day 12 was another grid based problem where each cell represented a plot with a plant that can be aggregated to a region
if they touch similar plots. Part was was simply counting the area and perimeter which was fairly straight forward.

Part 2 was somewhat frustrating because instead of the perimeter it was the sides of each region that needed to be 
counted. I did figure out that it was a corner counting thing but my mind simply glazed over thinking about each 
case, I ended up trying out claude to get the solution but I don't know if it was because of my code or my explanation
that were wrong but I could simply not prompt it to get the right solution. I ended up having to implement the corner
cases myself, it was only 8 conditions but I hate these kind of problems. I guess I should try better and think them
through. I wasted a bunch of time in the end on this because I kept prompting instead of stopping and thinking about
it myself. Once I did that, the solution was actually pretty straight forward.

I also learned that I had been abusing record classes for mutable data. The are primarily meant for immutable data
which is why the fields are final. Shows an area where Java is lacking where it can't make things fully immutable
when needed. Such as when I used a Map<> in a record, I could update it but not an int. It is definitely user error here,
I could also have read the docs more closely, but the language should be able to indicate that to me better too.
I ended up using plain classes instead but that showed me again how verbose Java is again.

# Day 13

The problem of today was essentially to figure out how many times mA + nB = X where A, B and X are given and we need
to find m and n. Part 1 was pretty easy to solve iteratively. Part 2 ended up needing a different solution since it was
essentially the same problem but with much larger numbers. Took me a while and I honestly didn't come up with the solution
myself. I did know that there must be a way to solve it mathematically, but I just don't remember how. I used Claude to 
get the formulas for how to do it. It was a bit of linear algebra.

# Day 14

Another grid and moving coordinates problem. Pretty simple. Part 2 was a bit vague, so I consulted the subreddit to 
understand what exactly I was supposed to expect as the output. 

# Day 15

Today was another 2d grid simulation problem. I still had a few problems with my initial solution. I should really
be getting better at those now. 

# Day 16

Day 16's problem was a path finding problem on a 2d grid. For part 1 I initially implemented BFS to find the min cost 
path to the end. But my initial implementation wasn't efficient. I ended up using Claude to get an A* implementation 
that worked well. 

Part 2 needed the insight of ditching A* in favor of Dijkstra's algorithm to find all paths that had min cost.
