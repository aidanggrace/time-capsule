"use client"
import { useEffect, useState } from "react"
import { Capsule } from "@/lib/types"
import { supabase } from "@/lib/supabase"

const API_URL = "http://localhost:8080"

export default function Home() {
  const [sentCapsules, setSentCapsules] = useState<Capsule[]>([])
  const [receivedCapsules, setReceivedCapsules] = useState<Capsule[]>([])
  const [username, setUsername] = useState("")
  const [showModal, setShowModal] = useState(false)
  const [recipientEmail, setRecipientEmail] = useState("")
  const [message, setMessage] = useState("")
  const [unlockAt, setUnlockAt] = useState("")

  const handleLogout = async () => {
    await supabase.auth.signOut()
    window.location.href = "/login"
  }

  const getAuthHeaders = async () => {
    const { data: { session } } = await supabase.auth.getSession()
    if (!session) {
      window.location.href = "/login"
      return null
    }
    return { Authorization: `Bearer ${session.access_token}` }
  }

  const fetchCapsules = async () => {
    const headers = await getAuthHeaders()
    if (!headers) return

    const [sentRes, receivedRes] = await Promise.all([
      fetch(`${API_URL}/capsules`, { headers }),
      fetch(`${API_URL}/received-capsules`, { headers }),
    ])

    const sent = await sentRes.json()
    const received = await receivedRes.json()
    setSentCapsules(sent || [])
    setReceivedCapsules(received || [])
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    const headers = await getAuthHeaders()
    if (!headers) return

    const res = await fetch(`${API_URL}/capsules`, {
      method: "POST",
      headers: { ...headers, "Content-Type": "application/json" },
      body: JSON.stringify({
        recipient_email: recipientEmail,
        message,
        unlock_at: unlockAt ? new Date(unlockAt).toISOString() : new Date().toISOString(),
      }),
    })

    if (res.ok) {
      setShowModal(false)
      setRecipientEmail("")
      setMessage("")
      setUnlockAt("")
      fetchCapsules()
    } else {
      alert("Failed to create capsule")
    }
  }

  useEffect(() => {
    const init = async () => {
      const { data: { session } } = await supabase.auth.getSession()
      if (session?.user.email) {
        setUsername(session.user.email.split("@")[0])
      }
      if (!session) {
        window.location.href = "/login"
        return
      }
      fetchCapsules()
    }
    init()
  }, [])

  return (
    <div>
      <div className="flex items-center justify-between">
        {username && <p>Welcome, {username}</p>}
        <div>
          <button onClick={() => setShowModal(true)} className="bg-blue-500 text-white px-3 py-1 rounded mr-2">
            Create Capsule
          </button>
          <button onClick={handleLogout} className="bg-red-500 text-white px-3 py-1 rounded">
            Logout
          </button>
        </div>
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center">
          <form onSubmit={handleCreate} className="bg-gray-900 text-white p-6 rounded flex flex-col gap-3 w-96">
            <h2 className="text-xl font-bold">Create Capsule</h2>
            <input
              type="email"
              placeholder="Recipient email"
              value={recipientEmail}
              onChange={(e) => setRecipientEmail(e.target.value)}
              className="border p-2 rounded text-black"
              required
            />
            <textarea
              placeholder="Your message"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              className="border p-2 rounded text-black"
              required
            />
            <label className="text-sm">Unlock date/time (leave empty for now)</label>
            <input
              type="datetime-local"
              value={unlockAt}
              onChange={(e) => setUnlockAt(e.target.value)}
              className="border p-2 rounded text-black"
            />
            <div className="flex gap-2">
              <button type="submit" className="bg-green-500 text-white px-3 py-1 rounded">Send</button>
              <button type="button" onClick={() => setShowModal(false)} className="bg-gray-600 text-white px-3 py-1 rounded">Cancel</button>
            </div>
          </form>
        </div>
      )}

      <h2 className="text-lg font-bold mt-4">Sent Capsules</h2>
      {sentCapsules.length > 0 ? (
        sentCapsules.map((c) => (
          <div key={c.id} className="border p-2 my-1">
            <p>To: {c.recipient_email}</p>
            <p>Message: {c.message}</p>
            <p>Unlocks: {c.unlock_at}</p>
          </div>
        ))
      ) : (
        <p>No sent capsules yet</p>
      )}

      <h2 className="text-lg font-bold mt-4">Received Capsules</h2>
      {receivedCapsules.length > 0 ? (
        receivedCapsules.map((c) => (
          <div key={c.id} className="border p-2 my-1">
            <p>From: {c.sender_email}</p>
            <p>Message: {c.message}</p>
            <p>Unlocks: {c.unlock_at}</p>
          </div>
        ))
      ) : (
        <p>No received capsules yet</p>
      )}
    </div>
  )
}

