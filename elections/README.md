# elections -- visualize US elections
![elections](elections.png)
![1864](1864.png)
![1964](1964.png)

## running

elections [options] file...

each file contains the layout and election results data as described below.
The repo contains results for all US presidential elections.

## interaction
* Left-arrow, Page-Down, Down-Arrow, Left-Mouse: move forward
* Right-arrow, Page-Up, Up-Arrow, Right-Mouse: move backward
* Home: first 
* End: last

## options
```
  -bgcolor string
        background color (default "black")
  -colsize float
        column size (canvas %) (default 7)
  -height int
        canvas height (default 900)
  -left float
        map left value (canvas %) (default 15)
  -rowsize float
        rowsize (canvas %) (default 9)
  -shape string
        shape for states:
        "c": circle,
        "h": hexagon,
        "s": square
        "l": line
        "p": plain text (default "c")
  -textcolor string
        text color (default "white")
  -textfont string
        font for text
  -top float
        map top value (canvas %) (default 75)
  -width int
        canvas width (default 1200)
```
## data

Tab-separated lists with fields: state, row, column, winner (r=republican, d=democrat, i=independent, f=Federalist, dr=Democratic-Republican, w=Whig), population.
The files begin with '# year canidate1 candidate2...

The party affiliation may be appended to the candidate name (for example, "Taylor:w", to indicate the Whig party).

For example:

```
# 1864 McClellan Lincoln Confederate
AL      6       6       i       964201
AR      5       4       i       435450
CA      4       0       r       379994
CT      3       9       r       460147
DE      4       9       d       112216
FL      7       8       i       140424
GA      6       7       i       1057286
IL      2       5       r       1711951
IN      3       5       r       1350428
IA      3       4       r       674913
KS      5       3       r       107206
KY      4       5       d       1155684
LA      6       4       r       708002
ME      0       10      r       628279
MD      4       8       r       687049
MA      2       9       r       1231066
MI      2       6       r       749113
MN      2       4       r       172023
MS      6       5       i       791305
MO      4       4       r       1182012
NV      3       1       r       6857
NH      1       10      r       326073
NJ      3       8       d       672035
NY      2       8       r       3880735
NC      5       7       i       992622
OH      3       6       r       2339511
OR      3       0       r       52465
PA      3       7       r       2906215
RI      3       10      r       174620
SC      5       6       i       703708
TN      5       5       r       1109801
TX      7       3       i       604215
VT      1       9       r       315098
VA      4       7       i       1219630
WV      4       6       r       376688
WI      1       5       r       775881
```