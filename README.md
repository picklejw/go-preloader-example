# Benchmarking Example and Results

Server-side rendering (SSR) has made a comeback as a strategy for faster page rendering, often promoted for performance and SEO benefits. While SSR can be useful, there are often more efficient ways to optimize page load times. Many websites already know the data needed for the user before making multiple network requests, yet current implementations often fetch this data in two or more sequential requests, creating unnecessary delays.

For example, when loading a product page, a typical workflow might be:

1. Request the HTML template for the product page.
2. Make a separate API call to fetch the product title, description, reviews, etc.

This staggered approach increases latency, server load, and hosting costs, particularly during traffic spikes.

## Alternative Approach

This demo shows an alternative approach to SSR—or a complementary method—that can reduce server load and improve perceived performance. By bundling the first API call directly with the `index.html` payload, the server can deliver essential data in a single request.

Benefits include:

* **Reduced latency**: Users see product details faster.
* **Lower server load**: Fewer requests reduce processing time.
* **Cost efficiency**: Handling multiple requests simultaneously costs more on servers or VPS instances.

In this benchmark, a Go-based optimized server was compared to a traditional staggered request setup (often implemented in Node.js). Go's efficiency allowed the optimized server to handle **2x–3x more sessions** under the same conditions.

> Note: "Sessions" here refer to network requests:
> A) Optimized: one request for `index.html` with API data bundled.
> B) Traditional: one request for `index.html`, then a separate API call.
> These benchmarks measure network requests only, without executing JavaScript in the browser.

## TODO
- The staggered request does not share the same TCP connection, these tests should be redone to reuse the same connection for better real world comparison.
- There should be a test case developed for comparing to using hydration techniques with a JavaScript runtime to compare actual benefits.

## Results

### Optimized Server

<img width="510" height="792" alt="optimized" src="https://github.com/user-attachments/assets/61ff7a28-45ff-42f5-8b4f-b5224c626948" />

### Traditional Staggered Requests

<img width="631" height="795" alt="staggered" src="https://github.com/user-attachments/assets/54c84ba9-1712-4b06-a101-bf45e038829c" />
