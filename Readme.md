# Project Index:
- [Main Doc](https://docs.google.com/document/d/1MnbWgbgAq2yvMtld9BMkyIBNHsPChyspgLmIX6H4eH0/edit?usp=sharing)
- [Meeting Notes](https://docs.google.com/document/d/1O7QY5yMG0Qxx3hBgZhlO9KaMPJBdpR8qYHzNmV9kHB0/edit?usp=sharing)
- [Trello Board](https://trello.com/b/yTZCJsdB/geometric-art-tool)

# Stuff to do
- [x] decide on the language/platform: *js. react native + expo*
- [ ] implement sketch of algorithm
- [ ] make a demo UI

## Backend

### Image Processing

Functions:
- [ ] Fitness function [delta e][1], [root mean square deviation][2]
- [ ] Average color

### Image Creation

Structures:
- [ ] Shapes (representation of triangles, rectangles, elipses, etc)
- [ ] Canvas with state. (The actual canvas to draw. To be able to resume from a state, or undo actions on the canvas.)

Functions:
- [ ] Generate random shape
- [ ] Append shape to canvas
- [ ] Export canvas to Image (SVG?)

### Variation of Hill Climbing

This variation borrows certain things from the genetic algorithms such as a fitness function and mutations.
Hill Climbing was chosen due to it's simplicity and efficiency[1]

Functions:
- [ ] Mutation (Vertex change of a triangle / radius change of an elipse)

Algorithm:
1. Generate Base (Average color of the image)
2. generate *n* random shapes
3. Create *n* new candidate states of the canvas
4. Rank by fitness these candidates and choose the first
5. Mutate the new state (only accept mutation if new image is better)
6. Repeat until desired complexity is achieved

## UI
*TODO*

---

**To Read!!!**  
[1]: http://zschuessler.github.io/DeltaE/learn/  
[2]: https://en.wikipedia.org/wiki/Root-mean-square_deviation  
[3]: https://sci-hub.st/10.1109/HICSS.1993.284069
