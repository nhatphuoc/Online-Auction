"use client"

import type React from "react"

import { useState } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { authService } from "@/services/auth-service"

export default function RegisterPage() {
  const [fullName, setFullName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [phoneNumber, setPhoneNumber] = useState("")
  const [loading, setLoading] = useState(false)
  const router = useRouter()

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    try {
      const res = await authService.register({ fullName, email, password, phoneNumber })
      if (res.success) {
        router.push(`/verify-otp?email=${encodeURIComponent(email)}`)
      } else {
        alert(res.message || "Registration failed")
      }
    } catch (err) {
      alert("An error occurred")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-black p-6">
      <Card className="w-full max-w-md bg-white/5 border-white/10 p-4">
        <CardHeader>
          <CardTitle className="text-3xl font-bold tracking-tighter text-center">Create Account</CardTitle>
          <p className="text-zinc-500 text-center text-sm">Join the world's most exclusive auction platform.</p>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <label className="text-xs font-bold uppercase tracking-widest text-zinc-500">Full Name</label>
              <Input
                value={fullName}
                onChange={(e) => setFullName(e.target.value)}
                required
                className="h-12 bg-white/5 border-white/10 rounded-xl"
              />
            </div>
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
            <div className="flex flex-col gap-2">
              <label className="text-xs font-bold uppercase tracking-widest text-zinc-500">Phone Number</label>
              <Input
                value={phoneNumber}
                onChange={(e) => setPhoneNumber(e.target.value)}
                className="h-12 bg-white/5 border-white/10 rounded-xl"
              />
            </div>
            <Button
              type="submit"
              disabled={loading}
              className="h-12 bg-white text-black hover:bg-zinc-200 rounded-full font-bold mt-2"
            >
              {loading ? "Creating..." : "Register"}
            </Button>
          </form>
          <p className="text-center mt-6 text-sm text-zinc-500">
            Already have an account?{" "}
            <Link href="/sign-in" className="text-white font-bold hover:underline">
              Sign In
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
