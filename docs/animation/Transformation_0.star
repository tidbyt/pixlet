animation.Transformation(
  child = render.Box(render.Circle(diameter = 6, color = "#0f0")),
  duration = 100,
  delay = 0,
  origin = animation.Origin(0.5, 0.5),
  direction = "alternate",
  fill_mode = "forwards",
  keyframes = [
    animation.Keyframe(
      percentage = 0.0,
      transforms = [animation.Rotate(0), animation.Translate(-10, 0), animation.Rotate(0)],
      curve = "ease_in_out",
    ),
    animation.Keyframe(
      percentage = 1.0,
      transforms = [animation.Rotate(360), animation.Translate(-10, 0), animation.Rotate(-360)],
    ),
  ],
),