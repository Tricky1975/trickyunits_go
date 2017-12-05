# GINI Is Not Ini 

GINI is nothing more but a quick configuration system I whipped up for 
easily configuring many of my programs.

It was originally set up in BlitzMax, parts of this Go version are
translation from BlitzMax to Go, however there are a few diffences, 
but those are no longer relevant.


~~~
[vars]
Hi=Hello
Earth=World

[list:Animals]
Dog
Cat
Elephant
Bird

[list:Plants,Flora]
Tree
Flower
Leaf

[rem]
Cool, huh?
~~~


There are basically three kinds of datablocks "vars", "list" and "rem".
Well "rem" is just for commenting the parser will ignore those.
"vars" is just for variable definition and "list" for making a list of strings.
And that is all there's to it. Simplicity was key for this. But since GINI has seen service in many of my apps it does show to be functional.

Ah yeah, the second list block in the example above, basically two list records "Plants" and "Flora" were created, but they share the same memory reference, so they'll always contain the same data. That feature may seem useless, but don't be fooled :P

A few notes
- GINI has only been setup to work with with the UNIX textfile format. So only "newline $a/LF" is used and not the "cursorreturn $d/CR". Windows users in particular will have to take note of this as this disqualifies notepad and a lot of standard Windows editors for making GINI files. Cross platform code editors (like jEdit, Geany and some others) are best, and look well if they are properly configured for the correct file format!)
- GINI's identifiers are CASE INSENSITIVE so in the example above "Hi" and "HI" and "hi" all contain the same value! 
- The same goes for tags so \[vars\] has the same effect as \[VARS\].



