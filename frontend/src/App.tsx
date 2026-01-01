import { useState, useCallback } from 'react';
import { Search, IngestItem } from "../wailsjs/go/main/App";

interface Item {
    id: number;
    name: string;
    description: string;
    type: string;
    icon: string;
    icon_large: string;
    members: boolean;
    current_price: number;
    current_trend: string;
    today_price_change: number;
    today_trend: string;
}

function App() {
    const [query, setQuery] = useState("");
    const [results, setResults] = useState<Item[]>([]);
    const [ingestId, setIngestId] = useState("4151");
    const [status, setStatus] = useState("");

    const handleSearch = useCallback((q: string) => {
        setQuery(q);
        if (q.length > 0) {
            Search(q).then((res: any) => {
                 setResults(res || []);
            });
        } else {
            setResults([]);
        }
    }, []);

    const handleIngest = () => {
        setStatus("Ingesting...");
        const id = parseInt(ingestId);
        if (isNaN(id)) {
            setStatus("Invalid ID");
            return;
        }
        IngestItem(id).then((msg) => {
            setStatus(msg);
            // Re-search if query is active
            if (query) handleSearch(query);
        });
    };

    return (
        <div id="App" style={{ padding: '20px', fontFamily: 'Nunito, sans-serif', maxWidth: '800px', margin: '0 auto' }}>
            <h1 style={{ textAlign: 'center', marginBottom: '30px' }}>RS Item Database</h1>
            
            <div style={{ marginBottom: '20px', padding: '15px', border: '1px solid #eee', borderRadius: '8px', background: '#fafafa' }}>
                <h3 style={{ marginTop: 0 }}>Debug / Ingest</h3>
                <div style={{ display: 'flex', gap: '10px' }}>
                    <input 
                        value={ingestId} 
                        onChange={(e) => setIngestId(e.target.value)} 
                        placeholder="Item ID (e.g. 4151)"
                        style={{ padding: '8px', borderRadius: '4px', border: '1px solid #ddd' }}
                    />
                    <button onClick={handleIngest} style={{ padding: '8px 16px', cursor: 'pointer', background: '#007bff', color: 'white', border: 'none', borderRadius: '4px' }}>Ingest Item</button>
                </div>
                <div style={{ marginTop: '10px', color: '#666', fontSize: '0.9em' }}>{status}</div>
            </div>

            <div className="search-box">
                <input 
                    className="input" 
                    value={query} 
                    onChange={(e) => handleSearch(e.target.value)} 
                    placeholder="Search for items..."
                    style={{ width: '100%', padding: '15px', fontSize: '18px', borderRadius: '8px', border: '1px solid #ccc', boxSizing: 'border-box' }}
                />
            </div>

            <div className="results" style={{ marginTop: '20px', textAlign: 'left' }}>
                {results.map((item) => (
                    <div key={item.id} style={{ display: 'flex', alignItems: 'center', marginBottom: '10px', background: 'white', padding: '15px', borderRadius: '8px', boxShadow: '0 2px 4px rgba(0,0,0,0.1)' }}>
                        <img src={item.icon} alt={item.name} style={{ marginRight: '20px', width: '40px', height: '40px' }} />
                        <div style={{ flex: 1 }}>
                            <div style={{ fontWeight: 'bold', fontSize: '1.2em' }}>{item.name}</div>
                            <div style={{ color: '#555', fontSize: '0.9em' }}>{item.description}</div>
                        </div>
                        <div style={{ textAlign: 'right' }}>
                            <div style={{ color: '#007bff', fontWeight: 'bold' }}>{item.current_price.toLocaleString()} gp</div>
                            <div style={{ fontSize: '0.8em', color: item.today_price_change >= 0 ? 'green' : 'red' }}>
                                {item.today_price_change > 0 ? '+' : ''}{item.today_price_change.toLocaleString()}
                            </div>
                        </div>
                    </div>
                ))}
                {results.length === 0 && query !== "" && <div style={{ textAlign: 'center', color: '#999', marginTop: '20px' }}>No results found.</div>}
            </div>
        </div>
    )
}

export default App