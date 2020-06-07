---
date: 2015-02-15T14:28:00+01:00
metadata:
  title: Extract color from image using JavaScript and Web Worker
  summary: A demonstration of how to extract the dominant color of an image with pure JavaScript. For performance optimization I utilized a web worker, which had a great impact.
  tags:
    - JavaScript
    - Web Worker
---

In the last week I was presented with the task to automatically extract a color from an
image in a webpage and use this as the background color. A quick research resulted in the
following steps:

1.  Create a `canvas` element and draw the image on its 2d context.

1.  Get the pixel data from the context using `getImageData` method.

    **Note:** When you want to load images from another domain, you need to add the `crossorigin`
    attribute to the `img` tag ([more information on the MDN](https://developer.mozilla.org/en-US/docs/Web/HTML/CORS_enabled_image)).

1.  Calulate the dominant pixel color.

Thankfully there is a nice JavaScript library out there, so we do not need to implement all
logic by outself: [Color Thief](http://lokeshdhakar.com/projects/color-thief/). The following
example extracts the color from an image and applies as the background color:

```html
<div>
	<img id="myImage" src="/images/example.png" />
</div>

<script>
	var myImage = document.getElementById("myImage");
	myImage.addEventListener("load", function () {
		var colorThief = new ColorThief(),
			color = colorThief.getColor(myImage);

		myImage.parentNode.style.backgroundColor =
			"rgb(" + color[0] + ", " + color[1] + ", " + color[2] + ")";
	});
</script>
```

This works well for a few small images. But when trying to do this for a couple of big images
at once, it slows down the browser and results in hangs/lags on the website. Wo don't want
that.

So I played a bit with [web workers](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/basic_usage)
to get the heavy color calculation off the main thread. The only thing, that needs to be done
on the main thread is creating the canvas and getting the pixel data from it. All the rest can
be done separatly.

To explain the use of web workers a little bit, imaging the image from the example above. In the
main JavaScript we have to follwing code in the image onload event:

```js
// Setting up the web worker.
var worker = new Worker("worker.js");
worker.addEventListener("message", function (e) {
	var color = e.data;

	myImage.parentNode.style.backgroundColor =
		"rgb(" + color[0] + ", " + color[1] + ", " + color[2] + ")";
});

// Starting the web worker.
// (Imagine the function getImageDataUsingCanvas to create a
// canvas, drawing the image on its context and then returning
// the image data.)
var imageData = getImageDataUsingCanvas(myImage);
worker.postMessage(imageData);
```

Then in the worker.js file, we have the following code:

```js
addEventListener("message", function (e) {
	var imageData = e.data;

	// The getColor function is a modified version of the
	// ColorThief.getColor function, which directly accepts
	// the pixel data as an argument.
	var color = getColor(imageData);

	postMessage(color);
});
```

I made up a full example with and without web workers, so you can see the impact it has on the
overall performance:

[https://embed.plnkr.co/s1KlIF/preview](https://embed.plnkr.co/s1KlIF/preview)
