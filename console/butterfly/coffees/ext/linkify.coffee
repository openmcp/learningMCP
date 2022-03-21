walk = (node, callback) ->
  for child in node.childNodes
    callback.call(child)
    walk child, callback

linkify = (text) ->
  # http://stackoverflow.com/questions/37684/how-to-replace-plain-urls-with-links
  urlPattern = (
    /\b(?:https?|ftp):\/\/[a-z0-9-+&@#\/%?=~_|!:,.;]*[a-z0-9-+&@#\/%=~_|]/gim)
  pseudoUrlPattern = /(^|[^\/])(www\.[\S]+(\b|$))/gim
  emailAddressPattern = /[\w.]+@[a-zA-Z_-]+?(?:\.[a-zA-Z]{2,6})+/gim
  text
    .replace(urlPattern, '<a href="$&">$&</a>')
    .replace(pseudoUrlPattern, '$1<a href="http://$2">$2</a>')
    .replace(emailAddressPattern, '<a href="mailto:$&">$&</a>')

tags =
  '&': '&amp;'
  '<': '&lt;'
  '>': '&gt;'

escape = (s) -> s.replace(/[&<>]/g, (tag) -> tags[tag] or tag)

Terminal.on 'change', (line) ->
  walk line, ->
    if @nodeType is 3
      val = @nodeValue
      linkified = linkify escape(val)
      if linkified isnt val
        newNode = document.createElement('span')
        newNode.innerHTML = linkified
        @parentElement.replaceChild newNode, @
        true
