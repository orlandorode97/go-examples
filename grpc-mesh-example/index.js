const resp = await fetch("http://136.114.195.162:8080/v1/ping");
const data = await resp.json();
console.log({ data });
