import { fetchWithAuth } from "@/lib/api-client"

export const authService = {
  async register(payload: any) {
    const res = await fetchWithAuth("/auth/register", {
      method: "POST",
      body: JSON.stringify(payload),
    })
    return res.json()
  },

  async verifyOtp(payload: { email: string; otpCode: string }) {
    const res = await fetchWithAuth("/auth/verify-otp", {
      method: "POST",
      body: JSON.stringify(payload),
    })
    return res.json()
  },

  async signIn(payload: any) {
    const res = await fetchWithAuth("/auth/sign-in", {
      method: "POST",
      body: JSON.stringify(payload),
    })
    return res.json()
  },

  async getProfile() {
    const res = await fetchWithAuth("/users/profile/me")
    return res.json()
  },
}
