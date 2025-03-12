
import open3d as o3d
import sys

basefile = sys.argv[1]

print(o3d.__version__ + "\n")

mesh = o3d.io.read_triangle_mesh("curve.stl")
o3d.io.write_triangle_mesh("curve.obj", mesh)
