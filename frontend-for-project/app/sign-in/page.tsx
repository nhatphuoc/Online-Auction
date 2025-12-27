"use client"

import type React from "react"

import { useState } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { authService } from "@/services/auth-service"

export default function SignInPage() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [loading, setLoading] = useState(false)
  const router = useRouter()

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    try {
      const res = await authService.signIn({ email, password })
      if (res.success) {
        localStorage.setItem("accessToken", res.accessToken)
        localStorage.setItem("refreshToken", res.refreshToken)
        router.push("/")
      } else {
        alert(res.message || "Invalid credentials")
      }
    } catch (err) {
      alert("Login failed")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-black p-6">
      <Card className="w-full max-w-md bg-white/5 border-white/10 p-4">
        <CardHeader>
          <CardTitle className="text-3xl font-bold tracking-tighter text-center">Welcome Back</CardTitle>
          <p className="text-zinc-500 text-center text-sm">Securely sign in to your auction account.</p>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <label className="text-xs font-bold uppercase tracking-widest text-zinc-500">Email Address</label>
              <Input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="h-12 bg-white/5 border-white/10 rounded-xl"
              />
            </div>
            <div className="flex flex-col gap-2">
              <label className="text-xs font-bold uppercase tracking-widest text-zinc-500">Password</label>
              <Input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="h-12 bg-white/5 border-white/10 rounded-xl"
              />
            </div>
            <Button
              type="submit"
              disabled={loading}
              className="h-12 bg-white text-black hover:bg-zinc-200 rounded-full font-bold mt-2"
            >
              {loading ? "Signing in..." : "Sign In"}
            </Button>
          </form>
          <div className="relative my-6">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-white/10"></div>
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-[#0c0c0c] px-2 text-zinc-500">Or continue with</span>
            </div>
          </div>
          <Button variant="outline" className="w-full h-12 rounded-full border-white/10 font-bold bg-transparent">
            Google
          </Button>
          <p className="text-center mt-6 text-sm text-zinc-500">
            Don't have an account?{" "}
            <Link href="/register" className="text-white font-bold hover:underline">
              Register
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
