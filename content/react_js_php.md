# JavaScript, React and PHP

The following is really meant as an educational piece for my wife. But enjoy if
you aren't my wife :).

## Intro

When you first visit a website you establish a connection between your browser
(the client) and the website's host server. This connection is handled through
what is called TLS (Transport Layer Security - HTTPS) and is maintained through
session cookies.

After the connection is established the server will process any server side code
such as PHP and then send the client all of the static assets (HTML, CSS, JS,
images, etc...) so that your browser can render them to you. After your browser
receives all of these assets it interprets them using the DOM (Document Object
Model). As part of this process JavaScript code will be interpreted using the
JavaScript Event Loop.

### Aside

A quick backstory of what is going on here. In web development we often talk
about the frontend and backend quite a bit. These terms are often mistaken for
gated content such as an admin dashboard (backend) and what the end user
sees (frontend). In reality the frontend is everything the browser has access to
and the backend is everything the server has access to. These will be things such
as the database, server side code such as PHP and more. The frontend is
everything that was sent to the client during the initial connection to the
server.

Browsers handle all of the interpretation of the frontend code and each browser
does it slightly different and with different technology, however, the
implementation is based on standards that have been set up by the World Wide Web
Consortium (W3C) for HTML and CSS, Ecma International through the TC39 committee
(Technical Committee 39) for JavaScript which maintains the ECMAScript standard.

Browsers implement "Engines" or "Interpreters" that follow the specifications
set by those standards in order to execute and render the code on the client.
You will hear terms like V8 or Gecko which are just names for these engines in
Chrome and Firefox respectively.

## Chapter 1

### Introduction to JavaScript

JavaScript is often confused with Java, less now than it used to be, but I still
hear it quite often. Briefly, JavaScript was a solution to a problem back in the
early days of the internet. The creation of it was contracted by Netscape (one
of the first browsers which eventually became firefox) in order to implement
scripting in the browser. It was named JavaScript purely for marketing as Java
was the most popular programming language at the time.

JavaScript was really a scripting language at the beginning and was never
intended for heavy programming. This means that the architectural decisions made
early on were void of potential future use cases such as memory management and
type safety. I won't get into those details unless you want to learn more but
just know that it has become something it was never intended to be. That is
partly why the ECMAScript standard was created.

That being said, the modern version of JavaScript is quite powerful and is
capable of doing pretty much anything despite it not necessarily being the best
or right tool for the job.

#### The Event Loop

JavaScript utilizes a concept known as the Event Loop to handle code execution.
At a high level you can think of this literally as an infinite loop running in
your browser with events getting loaded into a queue and then executed
sequentially. This Event Loop is initialized after all of the HTML and CSS has
been received and the JavaScript files have been received and loaded via HTML
`<script>` tags.

JavaScript is mainly used for manipulating the DOM (Document Object Model). The
DOM is a tree structure maintained by the browser that is essentially an
interpretation of the HTML. So for example:

```html
<div>
    <p>This is some content</p>
</div>
```

would be interpreted roughly as:

```bash
...
{
  type: 'element',
  tagName: 'div',
  children: [
    {
      type: 'element',
      tagName: 'p',
      children: [
        {
          type: 'text',
          content: 'This is some content'
        }
      ]
    }
  ]
}
```

Each element is represented as a "node" in the tree. JavaScript is able to
interact with each node by traversing the tree. You may have noticed that this
looks very familiar to JSON, and that is a good way of thinking about it.

Say we wanted to change the content within the `<p>` tag using JavaScript. We
could do that with the following code snippet:

```javascript
const pElement = document.querySelector('p');

pElement.innerText = "This was added by JavaScript";
```

Let's break down what is going on here. `const` is used to define a constant
which is a type of variable in JavaScript that, once declared, cannot be
changed also known as "immutable". To the right of the `=` is the data we are
storing in the declared `pElement` variable. `pElement` is just the name we gave
our variable. This makes it easier to access our data later.

`document.querySelector('p');` is almost self explanatory. We are using the
document object provided by the DOM which is the top level object. Each HTML
file is known as a document. We then query a selector which in this case is the
tag name. So now `pElement` stores the entire object data for our `p` tag.

There is a small issue with this code that doesn't really matter for our
example. The way we are selecting the `<p>` tag may inadvertently target the first
`<p>` tag on the page. In other words, if we had another one on the page we may be
changing the wrong one. This can be solved by using a more granular selector such as an ID or
Classname but it doesn't matter for this brief example.

`pElement.innerText` is a reference to the innerText data stored in the <p>
object. Prior to us changing it it would show the content we defined earlier
"This is some content". We are now changing that to "This was added by
JavaScript".

This is a very basic example of the power of JavaScript, the important thing to
keep in mind here is that this manipulation is done on the client side only and
only for the client that executed the code. This means that if you refresh your
browser it will be reset to the original code provided by the server. The server
doesn't even know what was executed and neither do other clients potentially
connected to the server.

If the server provides an API endpoint for interacting with the stored data we
could potentially execute that code in a way that would send the update to the
server, overwriting the content. This will be touched on later when
we discuss React.

Typically we wouldn't use JS for things like this, or at least very rarely. It
is more commonly used for adding "events" to elements. Let's go through an
example of adding a "click" event to a button.

```html
<button id="myBtn">Click Me</button>

<div class="container">Wow I am visible now!</div>
```
Let's set up some quick CSS...

```css
.container {
    display: none; /* Hide this div on initial page load (server side) */
}

.active {
    display: block; /* Create a class called active to change the display */
}
```

Now for the JS...

```javascript
const btn = document.querySelector("#myBtn"); //  grab the button by ID
const container = document.querySelector('.container'); // grab the div by class

btn.addEventListener("click", function() {
    container.classList.toggle = "active";
});
```

We start by creating our simple HTML structure. We then add some CSS to make
sure our `<div>` isn't visible on the page when it is first loaded.

Next we add an event listener to our `btn` which queues the click event into the
Event Loop. When a user clicks on this button the code inside the event listener
callback will be executed, in this case toggling whether or not the <div> has
the class "active" or not.

CSS is executing from top down hence Cascading Style Sheet, so the .active class
will override the .container class therefore changing the display from "none"
(hidden from the browser) to "block" (visible).

This example is a very common pattern in web development, albeit a small
example. In practice we end up having much more complex logic but even things as
simple as showing/hiding elements is a common practice say for making a Modal
popup on a page. You are essentially just toggling the display or adding a new
class to the element to add additional styles to it.

A key thing to note here is that the <div> is still visible in the DOM. That is
why we are able to query it in JavaScript. The `display: none;` just changes the
way the element is rendered in the window it does not effect how it appears in
the DOM.

CSS cannot effect the DOM, it can only change the way it is rendered. JavaScript
on the other hand CAN effect the DOM. So say we wanted to remove that <div> from
the DOM altogether. We can do that by changing the logic inside our "click"
event. (I will omit the previously defined code and just focus on this event)

```javascript
btn.addEventListener("click", () => {
    if (!container) return;
    container.remove();
});
```

In this example we first check to see if our container exists in the DOM. If it
doesn't we simply return (exit) out of the event. If it does exist we remove it. We
need to add that if statement because otherwise we will receive an error saying
that container doesn't exist on subsequent clicks. We call this a "Guard
Clause".

Unlike our previous example where we changed the visibility of the element using
CSS rules this new example actually manipulates the DOM and literally removes
the element from it. This is useful for say a todo list where you want to remove
an item from the list once it's complete. Again, this container will re-appear
when we refresh the page (well it will be invisible since we have the "display:
none" set, but it will be in the DOM).

We are starting to see a problem appear. What if we wanted to manipulate the DOM
in a way that is permanent? So say we remove our container and we want the
server to know that it was removed? This is where libraries like React can be
useful. We can still do this with pure (vanilla) JavaScript but it can be
trickier.

## Chapter 2

### PHP and the Server

In contrast to the DOM manipulation discussed in the previous section we
sometimes need to prevent certain things from ever even getting to the client.
This is where server side code such as php and runtimes such as nodejs come in.
We will skip nodejs in this section and just focus on php to define the
concepts.

A common example would be for user authentication. Our first line of defense is
at the web server level which adds authentication guards to specific page
routes, but sometimes we want to guard chunks of content rather than an entire
page.

PHP is great at this, among other things. PHP can be installed on the server and
executed directly by that machine. This is often used for making requests to a
database such as a relational database like MySQL. We can also use php to verify
if the client requesting access to a specific page or section of content should
be able to. This all happens on the server before any information gets sent back
to the client. This is in contrast to the examples we discussed earlier with
JavaScript. PHP does not allow for easy DOM manipulation either, not in the same
way that JavaScript does so doing things like showing and hiding an element
dynamically isn't feasible.

This is where server side JavaScript runtimes like nodejs come in. Similar to
php, nodejs can be installed on a server acting as a runtime interpreting code
being executed. This essentially converts JavaScript into a backend language
like php or python for example. This means we can potentially manipulate the DOM
on both the server side and client side. It opens up doors that allowed the web
to grow into where we are today, it also means we can use a unified language for
both the frontend and backend, and in some cases mobile with tools like React
Native.

### The rise of frontend frameworks

Server side code has been the standard since the beginning of the Internet. And
it remains a standard today, however, as computers have improved it has allowed
the rise of frontend frameworks such as React. Since the frontend is entirely
interpreted on the client it relies on the client having sufficient resources
(RAM, CPU, GPU (for some 3d technology)) as well as a fast internet connection
for handling multiple requests quickly.

Traditionally this was not efficient and it was faster to handle all logic on
the server which could have more resources. As the amount of users visiting
the site scaled so would the amount of server resources. Today, this problem
isn't as prevalent due to the rise of frontend frameworks. Websites can handle
the minimal amount of server side code such as establish connections, handling
user logins, form submissions and then allow the rest of the logic to be handled
by the client (browser). In some cases all of the logic is handled by the client
with technologies such as Firebase gaining popularity.

As we discussed in the previous sections we typically use JavaScript for client
side DOM manipulation paired with CSS to provide interactivity on a website. We
also explored the difference between nodejs and php and why you might choose one
vs the other. In this section we will dive into React and explore how nodejs has
turned JavaScript into a fullstack language.

I will keep this high level so I am going to skip the environment config. That
is explained in detail on Reacts website. What we will instead do is explore the
same example from earlier except using React, we will then discuss the benefits.

### React

Previously we showed an example of how JavaScript can manipulate the DOM. Below
is that same example using React.

```javascript
import {useState} from "react";

export default function ToggleComponent() {
    const [visible, setVisible] = useState(false)

    const handleClick = () => {
        setVisible(!visible);
    }

    return (
        <>
            <button onClick={handleClick}>Click Me</button>

            {visible && (
                <div>Wow I am visible now!</div>
            )}
        </>
    )
}
```

You will notice in this example some weird code going on here. React provides a
language syntax known as JSX. JSX allows us to write HTML directly in our
JavaScript code. Another thing you will notice is were are exporting a function,
this is known as a component. Typically these are written in Pascal case where
the first letters of each word are capitalized. You also probably noticed we are
importing something called useState from "react". This is where the magic really
happens and why react is incredibly powerful.

State is essentially the condition of a given item between component renders. A
render in React is whenever something changes such as a component reloading or
state changing. We can create containers to store state by the syntax seen in
the example `const [visible, setVisible] = useState(false)`. What we are seeing
here is two variables initialized as constants. The first one is the current
state which was initialized as "false", the second is a "setter" allowing us to
change the state. This is known as a getter and a setter in JS land.

We created a helper function called `handleClick` that  allows us to set the
state to the opposite of what it currently is, essentially creating a toggle.

In our return statement we are returning JSX (HTML with JavaScript Sugar) out of
the component. You will notice an attribute added to the button called
`onClick`. This is a special React attribute that allows us to bind our
JavaScript code to the HTML. In this case we are passing our `handleClick`
function in.

Beneath the button you will notice some more freaky looking syntax. The `{` or
curly brace is used to allow us to pass JavaScript variables or code into our
HTML. This is similar to how twig works in Laravel for example. We call these
concepts "Templating Languages".

`{visible && ()}` basically means return the code inside `()` if `visible` is
true. By default we set it to false, so when the component first loads the <div>
won't load. Clicking on our button will toggle the state of visible from true to
false.

Pretty clever huh? This is just a small example of what React can do. This can
be paired with API calls, or User authentication and much more to provide clean
and simple UIs.

Let's take a look at another example. Previously we manually changed the
innerText of our <p> tag. Let's provide a way to create a todo list and allow
the user to add todos and change the text on the fly. This is a very common
example demonstrating the power of react.

```javascript
import { useState } from 'react'

export default function TodoList() {
    const [todos, setTodos] = useState([])
    const [newTodo, setNewTodo] = useState('')

    const addTodo = () => {
        const nextId = todos.length ? Math.max(...todos.map(todo => todo.id)) + 1 : 1
        setTodos([...todos, { id: nextId, text: newTodo }])
        setNewTodo('')
    }

    const handleInputChange = (e, id) => {
        const newTodos = todos.map(todo => {
            if (todo.id === id) {
                return { ...todo, text: e.target.value }
            }
            return todo
        })
        setTodos(newTodos)
    }

    const removeTodo = id => {
        setTodos(todos.filter(todo => todo.id !== id))
    }

    return (
        <div>
            {todos.map(todo => (
                <div key={todo.id}>
                    <input
                        value={todo.text}
                        onChange={e => handleInputChange(e, todo.id)}
                    />
                    <button onClick={() => removeTodo(todo.id)}>Remove</button>
                </div>
            ))}
            <input
                value={newTodo}
                onChange={e => setNewTodo(e.target.value)}
                placeholder="Add new todo"
            />
            <button onClick={addTodo}>Add Todo</button>
        </div>
    )
}

```

Similar to our previous example, we are using useState here. We set our todos
state to an empty array initially. We are also creating a state variable called
newTodo. We can use this to track new todos to add to our todo list.

Our first function `addTodo` will check add new todos to the list of todos using
spread syntax `[...todos, ..]`. This basically means spread the current list
first then add the stuff after the comma. Next we reset the newTodo state to
empty.

In our JSX we are using `map` to iterate over the items in our todos state. This
allows us to write one chunk of HTML rather than have to manually enter the HTML
for each new item.

This example highlights powerful features of react that are are cumbersome to do
in JavaScript and impossible to do in server side code. React gets "built" and
then sent to the client as vanilla JavaScript. So to the end user it feels and
looks essentially the same as if we were to write it in pure JavaScript but it
makes the developer experience way smoother.
