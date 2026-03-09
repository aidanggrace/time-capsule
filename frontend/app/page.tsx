"use client"
import { useEffect, useState } from "react"
import { Capsule } from "@/lib/types"

const API_URL = "http://localhost:8080"

export default function Home() {
  const [capsules, setCapsules] = useState<Capsule[]>([])

  useEffect(() => {
    const getCapsules = async () => {
      const res = await fetch(`${API_URL}/capsules`)
      const data = await res.json()
      setCapsules(data)
    }
    getCapsules()
  }, [])

  return (
    <div>
      <h1>My Capsules</h1>
      {capsules.map((capsule) => (
        <div key={capsule.id} className="capsule">
          <p>Message: {capsule.message}</p>
          <p>Unlocks: {capsule.unlock_at}</p>
          <p>Sends to: {capsule.recipient_email}</p>
        </div>
      ))}
    </div>
  )
}

