# rt-go

rt-go is a small ray tracer written in Go.

Website:

```text
https://rt.aayushk.dev
```

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

For the browser viewer:

```sh
go run . viewer
```

Then open:

```text
http://localhost:8080
```

## Features

- CPU ray tracing
- Antialiasing
- Gamma correction
- Spheres, triangles, and triangle meshes
- OBJ loading
- MTL material loading
- Lambertian, metal, and dielectric materials
- BVH acceleration
- PPM image output
- Simple browser viewer

## Notes

Large OBJ files can take time to load because the renderer parses the mesh and builds a BVH before rendering.
