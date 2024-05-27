package server

// type mockRecorder struct{}
//
// func (m *mockRecorder) Record(ctx context.Context, key string, data *bytes.Buffer) error {
// 	return nil
// }
//
// func TestUpload(t *testing.T) {
// 	grpcServer := grpc.NewServer()
//
// 	videoServer := newVideoServer(&mockRecorder{})
// 	videoServer.register(grpcServer)
//
// 	lis, err := net.Listen("tcp", ":"+config.GetGrpcPort())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	go grpcServer.Serve(lis)
// 	defer grpcServer.Stop()
//
// 	conn, err := grpc.NewClient(":"+config.GetGrpcPort(), grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		t.Fatalf("failed to dial server: %v", err)
// 	}
// 	defer conn.Close()
//
// 	client := pb.NewVideoServiceWriterClient(conn)
//
// 	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
// 	defer cancel()
//
// 	stream, err := client.Upload(ctx)
// 	assert.NoError(t, err)
//
// 	req := &pb.UploadRequest{
// 		Key:  "key_1",
// 		Data: []byte{0, 1, 2, 3, 4, 5},
// 	}
//
// 	err = stream.Send(req)
// 	assert.NoError(t, err)
//
// 	err = stream.CloseSend()
// 	assert.NoError(t, err)
//
// 	res, err := stream.CloseAndRecv()
// 	assert.NoError(t, err)
//
// 	assert.Equal(t, pb.UploadStatus_UPLOAD_STATUS_SUCCESS, res.GetUploadStatus())
// }
