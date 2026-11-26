resource "descope_email_template" "big" {
  project_id = descope_project.test.id
  method     = "magiclink"
  name       = "big-template"
  subject    = "Sign in"
  html_body = join("\n", [
    "<html><body>",
    "<p>Hello {{.projectName}} — this body is deliberately larger than the inline threshold.</p>",
    "<p>literal dollar-brace: $${not_interpolated} and directive: %%{ nope }</p>",
    "<p>ünïcödé 🎯 and \"quotes\" and <a href=\"https://example.com?a=1&b=2\">links</a></p>",
    join("", [for i in range(40) : "<p>filler paragraph number ${i} with some more text to cross the threshold</p>"]),
    "</body></html>",
  ])
}
