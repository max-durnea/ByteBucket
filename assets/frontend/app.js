const apiBase = '' // same origin

function setOutput(selector, data) {
  document.querySelector(selector).textContent = JSON.stringify(data, null, 2)
}

function getTokens() {
  return {
    access: localStorage.getItem('bb_access'),
    refresh: localStorage.getItem('bb_refresh')
  }
}

function saveTokens(access, refresh) {
  if (access) localStorage.setItem('bb_access', access)
  if (refresh) localStorage.setItem('bb_refresh', refresh)
}

function clearTokens() {
  localStorage.removeItem('bb_access')
  localStorage.removeItem('bb_refresh')
}

async function postJson(path, body, withAuth = false) {
  const headers = {'Content-Type': 'application/json'}
  if (withAuth) {
    const t = getTokens().access
    if (t) headers['Authorization'] = 'Bearer ' + t
  }
  const res = await fetch(path, {
    method: 'POST',
    headers,
    body: JSON.stringify(body)
  })
  return res.json().catch(()=>({status: res.status}))
}

async function getJson(path, withAuth = false) {
  const headers = {}
  if (withAuth) {
    const t = getTokens().access
    if (t) headers['Authorization'] = 'Bearer ' + t
  }
  const res = await fetch(path, {headers})
  return res.json().catch(()=>({status: res.status}))
}

document.getElementById('registerForm').addEventListener('submit', async (e)=>{
  e.preventDefault()
  const f = e.target
  const body = {username: f.username.value, email: f.email.value, password: f.password.value}
  const data = await postJson('/api/users', body)
  setOutput('#authOutput', data)
})

document.getElementById('loginForm').addEventListener('submit', async (e)=>{
  e.preventDefault()
  const f = e.target
  const body = {email: f.email.value, password: f.password.value}
  const data = await postJson('/api/login', body)
  if (data.JWTtoken || data.jwt_token || data.jwtToken) {
    // some variations
    const token = data.JWTtoken || data.jwt_token || data.jwtToken
    saveTokens(token, data.RefreshToken || data.refresh_token || data.Refresh_Token)
  }
  if (data.jwt_token && data.refresh_token) {
    saveTokens(data.jwt_token, data.refresh_token)
  }
  setOutput('#authOutput', data)
})

document.getElementById('btnRefresh').addEventListener('click', async ()=>{
  const tokens = getTokens()
  if (!tokens.refresh) { setOutput('#authOutput', {error: 'no refresh token stored'}); return }
  const data = await postJson('/api/refresh', {refresh_token: tokens.refresh})
  if (data.access_token) {
    saveTokens(data.access_token, tokens.refresh)
  }
  setOutput('#authOutput', data)
})

document.getElementById('btnClearTokens').addEventListener('click', ()=>{ clearTokens(); setOutput('#authOutput', {ok:true}) })

document.getElementById('uploadForm').addEventListener('submit', async (e)=>{
  e.preventDefault()
  const fileInput = e.target.file
  if (!fileInput.files || fileInput.files.length === 0) return setOutput('#filesOutput', {error: 'no file selected'})
  const file = fileInput.files[0]
  // send the file to the server; server will store to /tmp and upload to S3
  const fd = new FormData()
  fd.append('file', file)
  try {
    const res = await fetch('/api/files', {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer ' + getTokens().access
      },
      body: fd
    })
    const data = await res.json().catch(()=>({status: res.status}))
    setOutput('#filesOutput', data)
  } catch (err) {
    console.error('Upload exception', err)
    setOutput('#filesOutput', {error: String(err)})
  }
})

document.getElementById('btnListFiles').addEventListener('click', async ()=>{
  const data = await getJson('/api/files', true)
  setOutput('#filesOutput', data)
  const list = document.getElementById('filesList')
  list.innerHTML = ''
  if (Array.isArray(data)) {
    data.forEach(f=>{
      const li = document.createElement('li')
      li.textContent = `${f.file_name} (${f.mime_type})`
      const btn = document.createElement('button')
      btn.textContent = 'Get URL'
      btn.onclick = async ()=>{
        const urlRes = await getJson('/api/files/' + f.id, true)
        setOutput('#filesOutput', urlRes)
        if (urlRes.url) {
          // Try to fetch the presigned URL and save as blob. If the fetch fails (usually CORS),
          // fall back to opening the URL in a new tab where direct access works.
          try {
            const r = await fetch(urlRes.url)
            if (!r.ok) {
              setOutput('#filesOutput', {status: r.status, ok: false})
              return
            }
            const blob = await r.blob()
            const link = document.createElement('a')
            const objectUrl = URL.createObjectURL(blob)
            link.href = objectUrl
            link.download = f.file_name || 'downloaded'
            document.body.appendChild(link)
            link.click()
            link.remove()
            URL.revokeObjectURL(objectUrl)
          } catch (err) {
            // Likely a CORS/network error; fallback to opening the URL in a new tab.
            console.warn('Fetch to presigned URL failed, falling back to opening URL:', err)
            setOutput('#filesOutput', {warning: 'failed to fetch presigned URL, opening in new tab', error: String(err)})
            try {
              window.open(urlRes.url, '_blank')
            } catch (openErr) {
              setOutput('#filesOutput', {error: 'Could not open URL: ' + String(openErr)})
            }
          }
        }
      }
      li.appendChild(btn)
      list.appendChild(li)
    })
  }
})

document.getElementById('btnReset').addEventListener('click', async ()=>{
  const res = await fetch('/admin/reset', {method: 'POST'})
  setOutput('#adminOutput', {status: res.status})
})

// show stored tokens on load
setOutput('#authOutput', {tokens: getTokens()})