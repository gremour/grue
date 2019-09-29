# UI features

- **Window** -- representation of underlying window. User don't have to deal with
window directly.
- **Surface** -- a layer on window having custom resolution (pixel size),
providing optional bufferized draw. Currently, if pixel size is set, the draw
is bufferized (using pixel.Canvas). There is always should be default (primary)
surface, but user can add more. Surface contains reference to window.
Window contains a list of created surfaces.
- **Widget** -- an element of UI. Can be of different type: panel, pushbutton, etc.
Widgets form hierarchical structure, having parent and children.
Widget belongs to a surface.

Note: Grue currently uses same coordinate system as OpenGL: 0, 0 is at left bottom of the window.
That is different from most UI libs, which have 0, 0 at left top of the window.
