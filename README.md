# rt-go

rt-go is a small path tracer written in Go. It renders spheres, boxes, quads,
triangles, OBJ meshes, matte/metal/glass materials, and emissive lights.

![Cornell render](outputs/connel.png)

![Demo render](outputs/demo.png)


## Usage

Place an OBJ model in:

```text
models/figurine.obj
```

If the model has a material file, place it in the same folder and make sure the OBJ points to it with `mtllib`.

Then run:

```sh
go run .
```

This writes:

```text
image.ppm
```

Edit `config.go` to change the camera, render quality, objects, materials, and
lights.

For the browser viewer:

```sh
go run . viewer
```

Then open:

```text
http://localhost:8080
```

## Features

- CPU path tracing
- Parallel scanline rendering
- Antialiasing
- Gamma correction
- Direct light sampling
- Weighted sampling targets
- One-sided diffuse lights
- Spheres, quads, boxes, triangles, and triangle meshes
- OBJ loading
- MTL material loading
- Smooth OBJ normals
- Lambertian, metal, dielectric, and diffuse light materials
- Translate and Y-axis rotation transforms
- BVH acceleration
- PPM image output with live scanline updates
- Simple browser viewer

## Notes

Large OBJ files can take time to load because the renderer parses the mesh and builds a BVH before rendering.

The default scene is defined in Go code for now.
