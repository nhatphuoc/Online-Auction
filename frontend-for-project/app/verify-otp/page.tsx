"use client"

import type React from "react"

import { useState, Suspense } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { authService } from "@/services/auth-service"

function VerifyOtpContent() {
  const [otpCode, setOtpCode] = useState("")
  const [loading, setLoading] = useState(false)
  const searchParams = useSearchParams()
  const router = useRouter()
  const email = searchParams.get("email") || ""

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    try {
      const res = await authService.verifyOtp({ email, otpCode })
      if (res.success) {
        alert("Verification successful! You can now sign in.")
        router.push("/sign-in")
      } else {
        alert(res.message || "Invalid OTP")
      }
    } catch (err) {
      alert("An error occurred")
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card className="w-full max-w-md bg-white/5 border-white/10 p-4">
      <CardHeader>
        <CardTitle className="text-3xl font-bold tracking-tighter text-center">Verify Identity</CardTitle>
        <p className="text-zinc-500 text-center text-sm">We've sent a code to {email}.</p>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="flex flex-col gap-6">
          <div className="flex flex-col gap-2">
            <label className="text-xs font-bold uppercase tracking-widest text-zinc-500 text-center">OTP Code</label>
            <Input
              value={otpCode}
              onChange={(e) => setOtpCode(e.target.value)}
              required
              className="h-14 text-center text-2xl tracking-[0.5em] font-bold bg-white/5 border-white/10 rounded-2xl"
              maxLength={6}
            />
          </div>
          <Button
            type="submit"
            disabled={loading}
            className="h-12 bg-white text-black hover:bg-zinc-200 rounded-full font-bold"
          >
            {loading ? "Verifying..." : "Confirm Code"}
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}

export default function VerifyOtpPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-black p-6">
      <Suspense fallback={null}>
        <VerifyOtpContent />
      </Suspense>
    </div>
  )
}
