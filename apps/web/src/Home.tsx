"use client";

import { useState } from 'react';
import { generateMagicLink } from '@auth/api-client';

function Home() {
    const [email, setEmail] = useState<string>("");

  return <div style={{
    display: 'flex',
    justifyContent:'center',
    alignItems:'center'
  }}>
    <h1>Welcome to react router dom and Query!</h1>
    <input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="Email" />
    <button
    onClick={() => generateMagicLink(email)}
    >
        Test
    </button>
  </div>;
}

export default Home;
  
