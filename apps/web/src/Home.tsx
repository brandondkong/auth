"use client";

import { useState } from 'react';
import { generateMagicLink, useLazyApi } from '@auth/api-client';

function Home() {
     const [email, setEmail] = useState<string>("");
    const { refetch, loading, response } = useLazyApi(() => generateMagicLink(email), [email]);

     return <div style={{
        display: 'flex',
        flexDirection: 'column',
        justifyContent:'center',
        alignItems:'center'
      }}>
        { loading ? "Loading" : undefined }
        <h1>Welcome to react router dom and Query!</h1>
        <input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="Email" />
        <button
        onClick={refetch}
        >
            Test
        </button>
        { JSON.stringify(response) }
      </div>;
}

export default Home;
  
