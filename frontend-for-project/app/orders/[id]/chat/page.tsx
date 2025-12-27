"use client"

import type React from "react"

import { useState, useEffect, useRef } from "react"
import { useParams } from "next/navigation"
import { Navbar } from "@/components/navbar"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card } from "@/components/ui/card"
import { fetchWithAuth } from "@/lib/api-client"
import { Send } from "lucide-react"

export default function OrderChatPage() {
  const { id } = useParams()
  const [messages, setMessages] = useState<any[]>([])
  const [input, setInput] = useState("")
  const [ws, setWs] = useState<WebSocket | null>(null)
  const scrollRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    async function setupChat() {
      try {
        // Step 1: Get WS info
        const wsInfoRes = await fetchWithAuth(`/order-websocket/${id}`)
        const wsInfo = await wsInfoRes.json()

        const token = localStorage.getItem("accessToken")
        const wsUrl = `${wsInfo.order_service_websocket_url}?orderId=${id}&X-User-Token=${encodeURIComponent(token!)}&X-Internal-JWT=${encodeURIComponent(wsInfo.internal_jwt)}`

        const socket = new WebSocket(wsUrl)
        socket.onmessage = (event) => {
          const data = JSON.parse(event.data)
          if (data.type === "message") {
            setMessages((prev) => [...prev, data.data])
          }
        }
        setWs(socket)

        // Step 2: Load history
        const historyRes = await fetchWithAuth(`/orders/${id}/messages`)
        const historyData = await historyRes.json()
        setMessages(historyData || [])
      } catch (err) {
        console.error("Failed to setup chat", err)
      }
    }
    setupChat()

    return () => {
      ws?.close()
    }
  }, [id])

  useEffect(() => {
    scrollRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [messages])

  const sendMessage = (e: React.FormEvent) => {
    e.preventDefault()
    if (!input.trim() || !ws) return

    ws.send(JSON.stringify({ type: "message", content: input }))
    setInput("")
  }

  return (
    <div className="h-screen flex flex-col bg-black">
      <Navbar />
      <main className="flex-1 container mx-auto px-6 pt-24 pb-6 flex flex-col gap-4 overflow-hidden">
        <div className="flex items-center justify-between border-b border-white/10 pb-4">
          <h1 className="text-2xl font-bold tracking-tighter">Order Chat #{id}</h1>
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
            <span className="text-xs font-bold text-zinc-500 uppercase tracking-widest">Connected</span>
          </div>
        </div>

        <Card className="flex-1 bg-white/5 border-white/10 flex flex-col overflow-hidden rounded-3xl">
          <div className="flex-1 overflow-y-auto p-6 flex flex-col gap-4">
            {messages.map((msg, i) => (
              <div key={i} className={`flex flex-col gap-1 max-w-[80%] ${msg.isMine ? "self-end" : "self-start"}`}>
                <div
                  className={`p-4 rounded-2xl text-sm ${
                    msg.isMine ? "bg-white text-black font-medium" : "bg-white/10 text-white"
                  }`}
                >
                  {msg.message}
                </div>
                <span className="text-[10px] text-zinc-500 px-2">
                  {new Date(msg.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}
                </span>
              </div>
            ))}
            <div ref={scrollRef} />
          </div>

          <form onSubmit={sendMessage} className="p-4 border-t border-white/10 bg-black/20">
            <div className="flex gap-2">
              <Input
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="Type your message..."
                className="h-12 bg-white/5 border-white/10 rounded-full px-6"
              />
              <Button type="submit" className="h-12 w-12 rounded-full bg-white text-black p-0 shrink-0">
                <Send className="w-5 h-5" />
              </Button>
            </div>
          </form>
        </Card>
      </main>
    </div>
  )
}
