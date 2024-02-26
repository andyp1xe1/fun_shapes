import { StatusBar } from 'expo-status-bar';
import { StyleSheet, Text, View } from 'react-native';
import * as ImagePicker from 'expo-image-picker';

import Button from './Button';

import ImageViewer from './ImageViewer';
const PlaceholderImage = require('./assets/images/main-image.png')


//TODO: https://docs.expo.dev/tutorial/image-picker/#use-the-selected-image
export default function App() {
  const pickImageAsync = async () => {
    let result = await ImagePicker.launchImageLibraryAsync({
      allowsEditing: true,
      quality: 1,
    });

    if (!result.canceled) {
      console.log(result);
    } else {
      alert('You did not select any image.');
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.heading}>
        Fun Shapes
      </Text>
      <ImageViewer placeholderImageSource={PlaceholderImage} />
      <Button label="Upload" onPress={pickImageAsync} />
      <Button label="Take a Photo" onPress={pickImageAsync} />
      <StatusBar style="auto"/>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#b79c8e',
    color: '#EBDBB2',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 25,
  },
  
  heading: {
    color: '#4d676e',
    fontSize: 36,
    paddingVertical: 16,
    fontFamily: 'monospace',
    fontWeight: 'bold',
  },
});
