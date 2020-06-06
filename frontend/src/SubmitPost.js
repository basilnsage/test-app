import React, { useState } from 'react';
import axios from 'axios';

export default () => {
    const [title, setTitle] = useState('');
    const [body, setBody] = useState('');
    const [author, setAuthor] = useState('');
    const onSubmit = async (event) => {
        event.preventDefault();
        await axios.post("http://localhost:8000/posts", {
            title,
            body,
            author,
        });
        setTitle("");
        setBody("");
        setAuthor("");
    }
    return (
        <div>
            <form onSubmit={onSubmit}>
                <div className="form-group">
                    <label>Title</label>
                    <input
                        value={title}
                        onChange={e => setTitle(e.target.value)}
                        className="form-control"
                    />
                    <label>Content</label>
                    <input
                        value={body}
                        onChange={e => setBody(e.target.value)}
                        className="form-control"
                    />
                    <label>Author</label>
                    <input
                        value={author}
                        onChange={e => setAuthor(e.target.value)}
                        className="form-control"
                    />
                </div>
                <button className="btn btn-primary">Submit</button>
            </form>
        </div>
    );
}