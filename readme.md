##LibScan -- readme.md


**Concise blurb for github release**

LibScan is like a special-purpose 'grep' for go libraries
------- about this long, to not get truncated in listing ----------


**Download and Contact Info**

Download: <https://github.com/PhilStephens/LibScan>  
[Written and tested for Windows only, using windows-specific package 'walk' as input]  
Gmail address: PhilipRStephens


**Preliminary Plan**

Current intent has more to do with summarizing the contents of a library (repository or package) than 
tracing the subset actually in use, aiming to be at least a little more useful than the 2 files of grep 
results augmented by impromptu additional greps that I have been using mainly to find out 'can library X 
do function Y, and with what syntax using what library files'.

At minimum, need to detect comments (both single line and multi-line) so can ignore them <done in the 
sense of ignoring them, via scanner pkg>; might also detect comments that fit criteria of godoc and cgo
<not attempted, and cannot do via scanner features, but might do otherwise>.

Main intent: detect declarations, eg const, of both single and multiple types; store that info, and 
format it for human readable print or file <nearly done>.  Full parsing of Go syntax is NOT planned.

More advanced <not yet started>: detect usage (eg w/i declaration of a func or method), store it in 
something like a database, format it for human readable print or file, and trace threads of dependancy.
This would go way beyond the minimal 'improvement on grep' goal, and is considered 'extra'.


**What it is**

<TBD>


**Selected Project Concepts**


**How to use**


**Wishlist items**  
(

**Selected Function & Struct Descriptions**



**Condensed History**  

The grep files that inspired this were created 2012.08.10, via 'grep -Inrw type * > ..\grep.type.all.txt'
& 'grep -Inrw func * > ..\grep.func.all.txt' in a shell window (Git Shell) w/i Windows XP.  Which misses
a number of declarations, such as 4 struct definitions w/i a grouped type declaration.  So I have had in
the back of my mind since then to do something more specific to the Go programming language; at minimum
slightly better than my grep files and at maximum a guide to what other parts of library are needed if
one uses one of its features.

Project started 2012.10.10 when finished previous practice project, WCSG_TravellingSalesmanProblem 
(goLibrarian might be of more use to others than WCSG when semi-mature, but everything I do currently 
in Go is a practice project).  Tentative plan is to share preliminary versions of this one as a work 
in progress, more than once a month but not daily.

