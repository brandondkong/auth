"use client";

import { useState } from 'react';
import { generateMagicLink, useLazyApi, verifyMagicLink } from '@auth/api-client';

function Home() {
    const [email, setEmail] = useState<string>("");
    const [token, setToken] = useState<string>("");
    const { refetch, loading, response } = useLazyApi(() => generateMagicLink(email), [email]);
    const { refetch: refetchToken, response: responseToken } = useLazyApi(() => verifyMagicLink(token), [token]);

     return <div style={{
        display: 'flex',
        flexDirection: 'column',
        justifyContent:'start',
        height: '100vh'
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
        {
            !loading ? <div>
            <input value={token} onChange={ (e) => setToken(e.target.value)} placeholder="Token" />
            <button onClick={refetchToken}>Verify Token</button>
            { JSON.stringify(responseToken) }
            </div>
            : undefined
        }
      </div>;
}

export default Home;
  
