document.addEventListener('DOMContentLoaded', function() {
  var input = document.querySelector('#input');
  var submit = document.querySelector('#submit');
  var output = document.querySelector('iframe#output');
  var outputWindow = output.contentWindow
  var outputDocument = outputWindow.document
  outputDocument.body.style.fontSize = window.getComputedStyle(input, null).fontSize;
  outputDocument.body.style.fontFamily = window.getComputedStyle(input, null).fontFamily;

  var onsubmit = function() {
    execute(outputWindow, outputDocument, input.value);
  };
  submit.addEventListener('click', onsubmit);
  onsubmit();

  document.querySelector("textarea").addEventListener("keydown", function(e) {
    if (e.keyCode === 9) {
      e.preventDefault();
      var target = e.target;
      var val = target.value;
      var pos = target.selectionStart;
      target.value = val.substr(0, pos) + '\t' + val.substr(pos, val.length);
      target.setSelectionRange(pos + 1, pos + 1);
    }
  });
});

function execute(window, document, code) {
  console.log(code);
  var out = document.querySelector('pre#out');
  if (!out) {
    out = document.createElement('pre');
    out.id = "out"

    out.style.fontSize = window.getComputedStyle(document.body, null).fontSize;
    out.style.fontFamily = window.getComputedStyle(document.body, null).fontFamily;
    document.body.appendChild(out);
  }
  out.textContent = "waiting for remote server...";

  window.console.log = function(a) {
    out.textContent += a + "\n";
  };
  window.console.error = function(a) {
    out.textContent += a + "\n";
  };
  var script = document.createElement("script")
  script.onload = function() {
    out.textContent = "";
  };
  script.src = "/sandbox.js?code=" + encodeURIComponent(code);
  document.body.appendChild(script);
}
