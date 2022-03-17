import axios from 'axios';
import React, { useState, useEffect } from 'react';
// import styled from 'styled-components';
import { Container, Row, Col } from 'react-bootstrap';
// import $ from 'jquery';

const id = "cluster04";
const url = "http://10.0.3.50:57575";
// const pw = "1234";


const Home = () => {
    const [data, setData] = useState(null);

    // const makeURL = (res) => {
    //     setData(res);
        
    // }

    const getData = async () => {
        try {
            const response = await axios.get(
                'http://10.0.3.50:3124/scheduling',
            );
            // await makeURL(response);
            setData(response.data);
            console.log(response.data);
        } catch(e) {
            console.log(e);
        }
    }

    const exec = e => {
        // console.log("try");
        var cmd = e + '\n';
        var iframe = document.getElementById('iframe_1');
        // var cmd = window.frames['iframe_1'].document
        // iframe.contentWindow.butterfly.send(e+'\n');
        // console.log(iframe.contentWindow);
        iframe.contentWindow.postMessage(cmd, url);
        console.log(cmd);
        // iframe.contentWindow.document
        // console.log(cmd);
    }

    useEffect(() => {
        // document.domain = "localhost";
        // setTimeout(pwd, 2000);

        fetch("http://localhost:57575").then(res => console.log(res))

        const head = () => {
            var method = "GET";
            var xhr = new XMLHttpRequest();

            xhr.open(method, url);

            xhr.setRequestHeader("Access-Control-Allow-Origin", "*");
            
        }

        const pwd = () => {
            exec(id);
            // iframe.focus()
        }
        head();
        pwd();
        getData();
        // console.log(data);
    }, [])

    return (
        <>
            <Container fluid className="mt-5">
                <Row className="text-center">
                    <Col sm={4} className="border">
                        <div>
                            <div>
                                <button>불러오기 김정호</button>
                            </div>
                            {data && <textarea rows={7} value={JSON.stringify(data, null, 2)} readOnly={true}/>}
                        </div>
                    </Col>
                    <Col sm={8} className="border" style={{width: "1000px", height: "1000px"}}>
                        <iframe src={url} title="mcp" id="iframe_1" style={{width: "100%", height: "100%"}}></iframe>
                    </Col>
                </Row>
            </Container>
        </>
    )
}

export default Home;
