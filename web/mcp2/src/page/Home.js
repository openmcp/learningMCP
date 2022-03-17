import axios from 'axios';
import React, { useState, useEffect } from 'react';
// import styled from 'styled-components';
import { Container, Row, Col } from 'react-bootstrap';

const id = "cluster04";
// const url
// const pw = "1234";


const Home = () => {
    const [data, setData] = useState(null);
    const getData = async () => {
        try {
            const response = await axios.get(
                'http://10.0.3.50:3124/scheduling',
            );
            setData(response.data);
            console.log(response.data);
        } catch(e) {
            console.log(e);
        }
    }
    const exec = e => {
        console.log("try");
        // var cmd = e;
        var iframe = document.getElementById('iframe_1');
        // var cmd = iframe.contentWindow
        //console.log(iframe)
        //iframe.contentWindow.butterfly.send(e+'\n');
        //iframe.contentWindow.out(e+'\n');
        //iframe.contentWindow.postMessage(e+'\n', 'http://keti.asuscomm.com:57575');
        // console.log(iframe);
        // console.log(cmd);
    }

    const pwd = () => {
        // var iframe = document.getElementById('iframe_1');
        exec(id);
        // exec(pw);
        // iframe.focus()
        console.log("success");
    }

    useEffect(() => {
        setTimeout(pwd, 2000);
        getData();
        // pwd();
    }, [])

    return (
        <>
            <Container fluid className="mt-5">
                <Row className="text-center">
                    <Col sm={4} className="border">
                        <div>
                            <div>
                                <button>불러오기</button>
                            </div>
                            {data && <textarea rows={7} value={JSON.stringify(data, null, 2)} readOnly={true}/>}
                        </div>
                    </Col>
                    <Col sm={8} className="border" style={{width: "1000px", height: "1000px"}}>
                        <iframe src="http://10.0.3.50:57575" title="mcp" id="iframe_1" style={{width: "100%", height: "100%"}}></iframe>
                    </Col>
                </Row>
            </Container>
        </>
    )
}

export default Home;