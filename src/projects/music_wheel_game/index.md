---
layout: layout-project.njk
tags: projects
metadata:
  title: Music Wheel Game
  short_description: An HTML5 music game.
  images:
    - name: main
      alt: Screenshot of gameplay
  sources:
    - type: try
      url: /demo/music-wheel-game
    - type: git
      url: https://github.com/frigus02/music-wheel-game
  tags:
    - Web Components
    - ES Modules
    - Web Audio
---

An HTML5 music game. I'm mostly building this to play around with some new stuff like:

- Web Components: [Custom Elements v1: Reusable Web Components](https://developers.google.com/web/fundamentals/getting-started/primers/customelements), [Shadow DOM v1: Self-Contained Web Components](https://developers.google.com/web/fundamentals/getting-started/primers/shadowdom)
- ES Modules: [ES6 Modules in Depth](https://ponyfoo.com/articles/es6-modules-in-depth), [ECMAScript modules in browsers](https://jakearchibald.com/2017/es-modules-in-browsers/)
- [Web Audio](https://developer.mozilla.org/en-US/docs/Web/API/Web_Audio_API)
- [Redux](https://redux.js.org/)

Gameplay:

- As the music plays, colored circles will be spawning in the middle.
- Collect them with the mouse to trigger their effect:
  - **blue**: increase points by the multiplicator
  - **yellow**: increase multiplicator (+ 1)
  - **red**: decrease multiplicator (/ 1.5)
  - **purple**: increase color circle size for a few seconds
- Try to get as much points a possible?! I haven't really figured this part out yet.
