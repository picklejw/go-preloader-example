let preloaderInstance

export default (baseApiPath) => {
  if (!preloaderInstance) {
    preloaderInstance = new PreloaderAPILoader(baseApiPath)
  } else {
    if (baseApiPath) {
      console.warn(`PreloaderAPI was already initilized with: '${preloaderInstance.apiBasePath}' and you tried to initilize with: '${baseApiPath}'`)
    }
  }
  return preloaderInstance
}

class PreloaderAPILoader {
  constructor(apiBasePath) {
    if (!apiBasePath) {
      throw new Error("PreloaderAPILoader requires the API base path for routes not preloaded in request.")
    }
    this.apiBasePath = apiBasePath
  }

  async get(path, options = {}) {
    // check preloaded cache
    if (typeof window !== "undefined" && window.httpPreload && window.httpPreload[path]) {
      const cached = window.httpPreload[path];
      try {
        return JSON.parse(cached.body); // return parsed body
      } catch {
        return cached.body; // fallback raw
      }
    }

    // fallback to real fetch
    const res = await fetch(`${this.apiBasePath}${path}`, { method: "GET", ...options });
    if (!res.ok) {
      throw new Error(`HTTP error ${res.status}`);
    }
    return res.json();
  }

  async post(path, data, options = {}) {
    // no cache for POST, always call
    const res = await fetch(`${this.apiBasePath}${path}`, {
      method: "POST",
      headers: { "Content-Type": "application/json", ...(options.headers || {}) },
      body: JSON.stringify(data),
      ...options
    });
    if (!res.ok) {
      throw new Error(`HTTP error ${res.status}`);
    }
    return res.json();
  }
}
