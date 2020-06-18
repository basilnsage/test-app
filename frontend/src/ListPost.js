import React, { useState, useEffect } from 'react';
import axios from 'axios';
import SubmitComment from "./SubmitComment";
import ListComment from "./ListComment";

export default () => {
    const [posts, setPosts] = useState({});
    const fetchPosts = async() => {
        const res = await axios.get("http://localhost:8002/posts");
        setPosts(res.data['posts']);
    };

    useEffect(() => {
        fetchPosts();
    }, []);

    const styledPosts = Object.values(posts).map(post => {
        console.log(posts)
        return (
            <div
                className="card"
                style={{width: "30%", marginBottom: "20px"}}
                key={post.id}
            >
                <div className="card-body">
                    <h3>{post.title}</h3>
                    <h4>{post.body}</h4>
                    <ListComment comments={post.comments} />
                    <SubmitComment postId={post.id} />
                    <h5>{post.createdAt}</h5>
                </div>
            </div>
        );
    });

    return <div className="d-flex flex-row flex-wrap justify-content-between">
        {styledPosts}
    </div>
}