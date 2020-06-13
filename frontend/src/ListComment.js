import React, {useState, useEffect} from 'react';
import axios from 'axios';

export default ({postId}) => {
    const [comments, setComments] = useState({});
    const fetchComments = async () => {
        const res = await axios.get(`http://localhost:8001/posts/${postId}/comments`)
        setComments(res.data['comments'])
    };

    useEffect(() => {
        fetchComments();
    }, []);

    const styledComments = Object.values(comments).map(comment => {
        return <li key={comment.id}>{comment.body}</li>;
    });
    return <ul>
        {styledComments}
    </ul>
};