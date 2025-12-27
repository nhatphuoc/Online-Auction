const API_BASE_URL = "http://localhost:8080/api"

export async function fetchWithAuth(endpoint: string, options: RequestInit = {}) {
  const token = typeof window !== "undefined" ? localStorage.getItem("accessToken") : null

  const headers = new Headers(options.headers)
  if (token) {
    headers.set("X-User-Token", token)
  }
  if (!headers.has("Content-Type") && !(options.body instanceof FormData)) {
    headers.set("Content-Type", "application/json")
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  })

  if (response.status === 401) {
    // Handle token expiration / refresh logic here
    if (typeof window !== "undefined") {
      localStorage.removeItem("accessToken")
      window.location.href = "/sign-in"
    }
  }

  return response
}
