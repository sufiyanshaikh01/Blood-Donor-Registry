import React, { useState, useEffect } from 'react'; 

const BloodDonorApp = () => {
  const [donors, setDonors] = useState([]);
  const [formData, setFormData] = useState({ name: '', blood_group: '', phone: '', city: '' });
  const [filter, setFilter] = useState('');
  
  // New States for UX
  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');

  const fetchDonors = async () => {
    setIsLoading(true);
    setErrorMessage(''); // Purane errors clear karein
    try {
      const url = filter ? `http://localhost:8080/donors?blood_group=${encodeURIComponent(filter)}` : 'http://localhost:8080/donors';
      const res = await fetch(url);
      
      if (!res.ok) throw new Error('Failed to fetch data from Go server');
      
      const data = await res.json();
      setDonors(data || []);
    } catch (err) {
      setErrorMessage('Server connection failed. Check if Go API is running.');
    } finally {
      setIsLoading(false);
    }
  };
const handleDelete = async (phone) => {
  try {
    // Check karein ki yahan http:// aur localhost:8080 sahi se likha hai
    const response = await fetch(`http://localhost:8080/donors?phone=${phone}`, { 
      method: 'DELETE' 
    });
    
    if (response.ok) {
      fetchDonors();
    } else {
      console.error("Server error");
    }
  } catch (error) {
    console.error("Backend Can't connect:", error);
  }
};

const handleEdit = (donor) => {
  setFormData(donor); 
  window.scrollTo(0, 0); 
};

  useEffect(() => { fetchDonors(); }, [filter]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const res = await fetch('http://localhost:8080/donors', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData),
      });
      if (res.ok) {
        setFormData({ name: '', blood_group: '', phone: '', city: '' });
        fetchDonors();
      }
    } catch (err) {
      alert("Error saving donor: " + err.message);
    }
  };

  return (
    <div style={styles.container}>
      <h2>🩸 Blood Donor Registry 🩸 </h2>
      
      {/* Registration Form */}
      <form onSubmit={handleSubmit} style={styles.form}>
        <input placeholder="Name" required value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} />
        <input placeholder="Blood Group (O+, A-, etc)" required value={formData.blood_group} onChange={e => setFormData({...formData, blood_group: e.target.value})} />
        <input placeholder="Phone" required value={formData.phone} onChange={e => setFormData({...formData, phone: e.target.value})} />
        <input placeholder="City" required value={formData.city} onChange={e => setFormData({...formData, city: e.target.value})} />
        <button type="submit" style={styles.button}>Register Donor</button>
      </form>

      <hr />

      {/* Filter and List */}
      <div style={styles.listHeader}>
        <h3>Donor List</h3>
        <select onChange={(e) => setFilter(e.target.value)} style={styles.select}>
          <option value="">Filter by Blood Group</option>
          <option value="O+">O+</option>
          <option value="O-">O-</option>
          <option value="A+">A+</option>
          <option value="A-">A-</option>
          <option value="B+">B+</option>
          <option value="B-">B-</option>
        </select>
      </div>

      {/* Conditional Rendering: Loading, Error, or List */}
      {isLoading && <div style={styles.loader}>⏳ Fetching donors from CSV...</div>}
      
      {errorMessage && <div style={styles.error}>{errorMessage}</div>}

      {!isLoading && !errorMessage && (
        <ul style={styles.list}>
          {donors.length > 0 ? donors.map((d, i) => (
            <li key={i} style={styles.listItem}>
  <div>
    <strong>{d.name}</strong> - <span style={{color: 'red'}}>{d.blood_group}</span> | 📞 {d.phone} | 📍 {d.city}
  </div>
  <div style={{ marginTop: '10px' }}>
    <button onClick={() => handleEdit(d)} style={styles.editBtn}>
      ✏️ Edit
    </button>
    <button onClick={() => handleDelete(d.phone)} style={styles.deleteBtn}>
      🗑️ Delete
    </button>
  </div>
</li>
          )) : <p>No donors found.</p>}
        </ul>
      )}
    </div>
  );
};

// Basic Styles
const styles = {
  container: { maxWidth: '600px', margin: 'auto', padding: '20px', fontFamily: 'sans-serif' },
  form: { display: 'flex', flexDirection: 'column', gap: '10px', marginBottom: '20px' },
  button: { padding: '10px', background: '#d32f2f', color: 'white', border: 'none', cursor: 'pointer' },
  listHeader: { display: 'flex', justifyContent: 'space-between', alignItems: 'center' },
  select: { padding: '5px', height: '30px' },
  loader: { textAlign: 'center', padding: '20px', color: '#666' },
  error: { color: 'red', background: '#fee', padding: '10px', borderRadius: '5px', marginBottom: '10px' },
  list: { listStyle: 'none', padding: 0 },
  listItem: { padding: '10px', borderBottom: '1px solid #ddd' }
};

export default BloodDonorApp;