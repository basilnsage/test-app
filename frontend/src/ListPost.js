import React, { useState, useEffect } from 'react';
import axios from 'axios';
import SubmitComment from "./SubmitComment";

export default () => {
    const [posts, setPosts] = useState({});
    const fetchPosts = async() => {
        const res = await axios.get("http://localhost:8000/posts");
        setPosts(res.data['posts']);
    };

    useEffect(() => {
        fetchPosts();
    }, []);

    console.log(posts)

    const styledPosts = Object.values(posts).map(post => {
        return (
            <div
                className="card"
                style={{width: "30%", marginBottom: "20px"}}
                key={post.ID}
            >
                <div className="card-body">
                    <h3>{post.title}</h3>
                    <h4>{post.body}</h4>
                    <h4>{post.author}</h4>
                    <SubmitComment postId={post.ID} />
                    <h5>{post.createdAt}</h5>
                </div>
            </div>
        );
    });

    return <div className="d-flex flex-row flex-wrap justify-content-between">
        {styledPosts}
    </div>
}