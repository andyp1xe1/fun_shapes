# Stuff to do
- [ ] decide on the language
- [ ] implement sketch of algorithm
- [ ] make a demo UI

## Backend

### Image Processing

Functions:
- [ ] Fitness function ([Delta E](http://zschuessler.github.io/DeltaE/learn/) or [Root-mean-square deviation](https://en.wikipedia.org/wiki/Root-mean-square_deviation))
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

Delta E: http://zschuessler.github.io/DeltaE/learn/
Root-mean-square deviation: https://en.wikipedia.org/wiki/Root-mean-square_deviation
Python library: https://colour.readthedocs.io/en/latest/index.html
