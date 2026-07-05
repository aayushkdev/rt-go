# JSON scene guide

Scenes are JSON files. Start from `examples/gallery.json` or
`examples/cornell.json`, then edit the values.

## 1. File shape

Every scene has four top-level fields:

```json
{
  "output": "image.ppm",
  "camera": {},
  "render": {},
  "objects": []
}
```

Field | Required | Meaning
--- | --- | ---
`output` | yes | File path for the rendered image.
`camera` | yes | Camera position, framing, and focus.
`render` | yes | Quality, sampling, and worker settings.
`objects` | yes | List of objects in the scene.

## 2. Camera

The camera decides where the image is taken from and what it looks at.

```json
"camera": {
  "size": 600,
  "aspect": 1,
  "fov": 40,
  "from": [0, 1.5, 4],
  "look_at": [0, 1.5, -2],
  "focus": 6,
  "defocus": 0
}
```

Field | Required | Meaning
--- | --- | ---
`size` | yes | Image width in pixels. Height is calculated from `aspect`.
`aspect` | yes | Width divided by height. Use `1` for square, `1.777` for 16:9.
`fov` | yes | Vertical field of view in degrees. Higher values zoom out.
`from` | yes | Camera position `[x, y, z]`.
`look_at` | yes | Point the camera looks at `[x, y, z]`.
`focus` | yes | Focus distance for depth of field.
`defocus` | no | Lens blur angle. Use `0` or omit it for no blur.

## 3. Render settings

Render settings control quality and speed.

```json
"render": {
  "samples": 100,
  "max_depth": 20,
  "workers": 0,
  "flush_every_scanline": 10,
  "background": [0, 0, 0],
  "sky": false,
  "sampling_target_weight": 0.5
}
```

Field | Required | Meaning
--- | --- | ---
`samples` | yes | Rays per pixel. Higher values reduce noise but render slower.
`max_depth` | yes | Maximum ray bounces. Higher values allow more indirect light.
`workers` | no | Number of CPU workers. Use `0` to use all CPUs.
`flush_every_scanline` | no | How often partial PPM output is written while rendering.
`background` | yes | Background color `[r, g, b]` when `sky` is false.
`sky` | no | Use a blue sky gradient background instead of `background`.
`sampling_target_weight` | no | Chance to sample lights/targets instead of only material scattering. `0.5` is a good default.

Normal colors are usually between `0` and `1`. Light colors can be much higher,
like `[15, 15, 15]`.

## 4. Materials

Objects use a `material` field unless noted otherwise.

### 4.1 Matte

Diffuse, rough material.

```json
"material": {
  "type": "matte",
  "color": [0.7, 0.7, 0.7]
}
```

### 4.2 Metal

Reflective material.

```json
"material": {
  "type": "metal",
  "color": [0.8, 0.8, 0.75],
  "fuzz": 0.05
}
```

`fuzz` controls reflection blur. `0` is mirror-like.

### 4.3 Glass

Transparent refractive material.

```json
"material": {
  "type": "glass",
  "refraction": 1.5
}
```

`1.5` is a normal glass-like refraction value.

### 4.4 Light

Emissive material.

```json
"material": {
  "type": "light",
  "color": [10, 10, 10]
}
```

Field | Used by | Meaning
--- | --- | ---
`type` | all | `matte`, `metal`, `glass`, or `light`.
`color` | matte, metal, light | Base color or emission color.
`fuzz` | metal | Reflection roughness.
`refraction` | glass | Refraction index.

## 5. Common object fields

These fields can be used on most objects.

Field | Meaning
--- | ---
`material` | Object material.
`move` | Move object after rotation: `[x, y, z]`.
`rotate_x` | Rotate around the X axis in degrees.
`rotate_y` | Rotate around the Y axis in degrees.
`rotate_z` | Rotate around the Z axis in degrees.
`light` | Mark the object as an emissive light and sampling target.
`sample` | Force whether this object is sampled directly.

Rotation order:

```text
rotate_x -> rotate_y -> rotate_z -> move
```

Glass and light objects are sampled by default. Use `sample: true` for another
important object that should be sampled directly.

## 6. Object types

Each object in `objects` needs a `type`.

### 6.1 Sphere

```json
{
  "type": "sphere",
  "center": [0, 0.5, -2],
  "radius": 0.5,
  "material": { "type": "matte", "color": [0.8, 0.2, 0.2] }
}
```

Field | Required | Meaning
--- | --- | ---
`center` | yes | Sphere center `[x, y, z]`.
`radius` | yes | Sphere radius. Must be greater than zero.
`material` | yes | Sphere material.

### 6.2 Box

```json
{
  "type": "box",
  "min": [-0.5, 0, -0.5],
  "max": [0.5, 1, 0.5],
  "rotate_y": 20,
  "move": [0, 0, -2],
  "material": { "type": "matte", "color": [0.7, 0.7, 0.7] }
}
```

Field | Required | Meaning
--- | --- | ---
`min` | yes | Minimum corner `[x, y, z]`.
`max` | yes | Maximum corner `[x, y, z]`.
`material` | yes | Box material.

### 6.3 Triangle

```json
{
  "type": "triangle",
  "a": [-0.5, 0, -2],
  "b": [0.5, 0, -2],
  "c": [0, 1, -2],
  "material": { "type": "matte", "color": [0.2, 0.6, 0.9] }
}
```

Field | Required | Meaning
--- | --- | ---
`a` | yes | First point `[x, y, z]`.
`b` | yes | Second point `[x, y, z]`.
`c` | yes | Third point `[x, y, z]`.
`material` | yes | Triangle material.

### 6.4 Quad

```json
{
  "type": "quad",
  "q": [-1, 0, -2],
  "u": [2, 0, 0],
  "v": [0, 1, 0],
  "material": { "type": "matte", "color": [0.7, 0.7, 0.7] }
}
```

`q` is one corner. `u` and `v` are the two edge vectors.

Field | Required | Meaning
--- | --- | ---
`q` | yes | Starting corner `[x, y, z]`.
`u` | yes | First edge vector.
`v` | yes | Second edge vector.
`material` | yes | Quad material.

### 6.5 Model

```json
{
  "type": "model",
  "path": "models/figurine.obj",
  "height": 1.2,
  "position": [0, 0, -2]
}
```

Field | Required | Meaning
--- | --- | ---
`path` | yes | OBJ file path.
`height` | no | Scale model to this height.
`position` | no | Place the model on the ground at `[x, y, z]`.

Models use material data from the OBJ/MTL files. A JSON `material` field on a
model object is ignored.

The model loader recenters the model on X/Z and places its lowest point at
`position.y`. After that, `rotate_x`, `rotate_y`, `rotate_z`, and `move` still
apply like other objects.

### 6.6 Floor helper

```json
{
  "type": "floor",
  "x1": -2,
  "z1": -4,
  "x2": 2,
  "z2": 0,
  "y": 0,
  "material": { "type": "matte", "color": [0.7, 0.7, 0.7] }
}
```

Field | Required | Meaning
--- | --- | ---
`x1`, `z1` | yes | First floor corner.
`x2`, `z2` | yes | Opposite floor corner.
`y` | yes | Floor height.
`material` | yes | Floor material.

### 6.7 Wall on X

```json
{
  "type": "wall_x",
  "x": -2,
  "y1": 0,
  "y2": 3,
  "z1": 0,
  "z2": -4,
  "material": { "type": "matte", "color": [0.7, 0.7, 0.7] }
}
```

Field | Required | Meaning
--- | --- | ---
`x` | yes | Fixed X position.
`y1`, `y2` | yes | Bottom and top height.
`z1`, `z2` | yes | Wall depth range.
`material` | yes | Wall material.

### 6.8 Wall on Z

```json
{
  "type": "wall_z",
  "z": -4,
  "x1": -2,
  "x2": 2,
  "y1": 0,
  "y2": 3,
  "material": { "type": "matte", "color": [0.7, 0.7, 0.7] }
}
```

Field | Required | Meaning
--- | --- | ---
`z` | yes | Fixed Z position.
`x1`, `x2` | yes | Wall width range.
`y1`, `y2` | yes | Bottom and top height.
`material` | yes | Wall material.

### 6.9 Ceiling light helper

```json
{
  "type": "ceiling_light",
  "x1": -0.5,
  "z1": -2.5,
  "x2": 0.5,
  "z2": -1.5,
  "y": 3,
  "color": [15, 15, 15]
}
```

`ceiling_light` creates an emissive quad and automatically marks it as a light
and sampling target.

Field | Required | Meaning
--- | --- | ---
`x1`, `z1` | yes | First light corner.
`x2`, `z2` | yes | Opposite light corner.
`y` | yes | Light height.
`color` | yes | Light emission color. Higher values make the light brighter.

## 7. OBJ and MTL support

OBJ model loading is used by `type: "model"`.

Supported OBJ lines:

- `v x y z` vertices
- `vn x y z` vertex normals for smooth shading
- `f ...` faces
- `mtllib file.mtl` material library references
- `usemtl name` material assignment

Supported face forms:

```text
f 1 2 3
f 1/1 2/2 3/3
f 1//1 2//2 3//3
f 1/1/1 2/2/2 3/3/3
```

Notes:

- Faces with more than three vertices are triangulated.
- Positive and negative OBJ indices are supported.
- `mtllib` paths are resolved relative to the OBJ file.
- If a face has `vn` normals, the mesh uses smooth normals.
- If a face has no `vn` normals, it uses flat triangle normals.
- Texture coordinates and image textures are ignored for now.

MTL mapping:

MTL field | How rt-go uses it
--- | ---
`newmtl` | Starts a named material.
`Kd r g b` | Diffuse/base color.
`Ks r g b` | Specular color. Strong `Ks` makes the material metal.
`Ns value` | Shininess. Higher `Ns` makes metal less fuzzy.
`Ni value` | Refraction index for glass.
`d value` | Alpha. Values below `0.95` become glass.
`Tr value` | Transparency. Converted to alpha with `1 - Tr`.
`illum value` | Helps classify glass and metal.

Mapping rules:

- Transparent MTL materials become `glass`.
- Specular MTL materials become `metal`.
- Everything else becomes `matte`.
- Missing MTL files or unknown `usemtl` names fall back to mirror-like metal.
- Texture maps like `map_Kd` are ignored.

## 8. Quick object template

Use this when adding a new object:

```json
{
  "type": "sphere",
  "center": [0, 0.5, -2],
  "radius": 0.5,
  "rotate_x": 0,
  "rotate_y": 0,
  "rotate_z": 0,
  "move": [0, 0, 0],
  "sample": false,
  "material": { "type": "matte", "color": [0.7, 0.7, 0.7] }
}
```
